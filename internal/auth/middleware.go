package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/shared/response"
	"github.com/golang-jwt/jwt/v5"
)

func (a *Auth) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		unauthorizedErr := response.New().Error("Protected endpoint. You need credentials to continue", "00005")

		tokenString, ok := bearerToken(r)
		if !ok {
			unauthorizedErr.Send(w, http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			return a.secret, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) && token != nil && a.DB != nil {
				if sid, ok := claims["sid"].(string); ok {
					a.DB.DeleteSession(r.Context(), sid)
				}
			}
			unauthorizedErr.Send(w, http.StatusUnauthorized)
			return
		}

		user, err := userFromClaims(claims)
		if err != nil {
			unauthorizedErr.Send(w, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func bearerToken(r *http.Request) (string, bool) {
	h := r.Header.Get("Authorization")
	token, ok := strings.CutPrefix(h, "Bearer ")
	return strings.TrimSpace(token), ok
}