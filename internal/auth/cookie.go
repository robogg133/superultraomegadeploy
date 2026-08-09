package auth

import (
	"net/http"
	"time"
)

const (
	PersistCookie     = "persist_session"
	PersistCookiePath = "/api/v1/auth/refresh"
	persistTTL        = 30 * 24 * time.Hour
)

// ponytail: Secure=false for local development, turn on TLS.

func setPersistCookie(w http.ResponseWriter, sid string) {
	http.SetCookie(w, &http.Cookie{
		Name:     PersistCookie,
		Value:    sid,
		Path:     PersistCookiePath,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(persistTTL.Seconds()),
	})
}