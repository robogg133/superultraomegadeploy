package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestMain(m *testing.M) {
	os.Setenv("DEBUG_ENABLED", "true")
	os.Exit(m.Run())
}

func TestChiLoggerEmitsRequestLine(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf)

	handler := requestIDMiddleware(ChiLogger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/some/path", nil))

	out := buf.String()
	t.Log(out)
	if !strings.Contains(out, "Request") || !strings.Contains(out, "rid") || !strings.Contains(out, "/some/path") {
		t.Fatalf("missing request log, got: %q", out)
	}
}
