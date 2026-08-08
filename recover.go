package main

import (
	"net/http"
	"runtime/debug"

	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/shared/response"
	"github.com/go-chi/chi/v5/middleware"
)

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}

				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				}

				if r.Header.Get("Connection") != "Upgrade" {
					response.New().Error("Internal Server Error", "unknown").Send(w, 500)
				}
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
