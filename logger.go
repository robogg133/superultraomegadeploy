package main

import (
	"net/http"
	"os"
	"time"

	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/auth"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func initLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if DebugEnabled {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		log.Info().Msg("Debug enabled!")
		return
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

type logEntry struct {
	req *http.Request
}

func ChiLogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		entry := NewLogEntry(r)
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		t1 := time.Now()
		defer func() {
			entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1), nil)
		}()
		req := middleware.WithLogEntry(r, entry)
		req = req.WithContext(log.With().Str("rid", req.Context().Value(CTXKeyRequestID).(string)).Logger().WithContext(req.Context()))

		entry.req = req
		next.ServeHTTP(ww, req)
	}
	return http.HandlerFunc(fn)
}

func NewLogEntry(r *http.Request) *logEntry {
	return &logEntry{
		req: r,
	}
}

func (l *logEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, _ any) {
	if l.req.URL.Path == "/healthz" {
		return
	}
	lg := log.Ctx(l.req.Context()).Info().
		Str("method", l.req.Method).
		Str("ip", middleware.GetClientIPAddr(l.req.Context()).String()).
		Str("path", l.req.URL.Path).
		Int("status", status).
		Int("bytes", bytes).
		Str("elapsed", elapsed.String())
	uid, ok := auth.UserID(l.req.Context())
	if ok {
		lg = lg.Str("uid", uid.String())
	}
	lg.Msg("Request")

}

func (l *logEntry) Panic(v any, stack []byte) {
	log.Log().
		Str("level", "panic").
		Any("value", v).
		Str("stack", string(stack)).
		Send()
}
