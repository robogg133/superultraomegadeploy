package main

import (
	"context"
	"net/http"

	"github.com/rs/xid"
)

const CTXKeyRequestID int = 0

func requestIDMiddleware(next http.Handler) http.Handler {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(
			r.Context(),
			CTXKeyRequestID,
			xid.New().String(),
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
	return http.HandlerFunc(fn)
}
