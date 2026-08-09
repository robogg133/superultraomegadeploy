package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestRequireAuth(t *testing.T) {
	secret := "test-secret"
	a := New(nil, secret)

	sign := func(s string) string {
		tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":        "66f1b5b5-5b5b-4b5b-8b5b-5b5b5b5b5b5b",
			"sid":        "test-session",
			"email":      "a@b.c",
			"first_name": "Foo",
			"last_name":  "Bar",
			"exp":        time.Now().Add(time.Hour).Unix(),
		})
		s, err := tok.SignedString([]byte(s))
		if err != nil {
			t.Fatal(err)
		}
		return s
	}

	valid := func(w http.ResponseWriter, r *http.Request) {
		u, ok := ContextUser(r.Context())
		if !ok || u.Email != "a@b.c" || u.FirstName != "Foo" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
	handler := a.RequireAuth(http.HandlerFunc(valid))

	authHeader := func(token string) *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("no header", func(t *testing.T) {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("got %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("valid token", func(t *testing.T) {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, authHeader(sign(secret)))
		if rec.Code != http.StatusNoContent {
			t.Fatalf("got %d, want %d", rec.Code, http.StatusNoContent)
		}
	})

	t.Run("wrong secret", func(t *testing.T) {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, authHeader(sign("other")))
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("got %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("expired", func(t *testing.T) {
		tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": "66f1b5b5-5b5b-4b5b-8b5b-5b5b5b5b5b5b",
			"sid": "test-session",
			"exp": time.Now().Add(-time.Hour).Unix(),
		})
		s, err := tok.SignedString([]byte(secret))
		if err != nil {
			t.Fatal(err)
		}
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, authHeader(s))
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("got %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})
}