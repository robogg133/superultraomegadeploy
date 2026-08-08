package main

import (
	"net/http"

	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/shared/response"
	"git.servidordomal.lol/robogg133/superultraomegadeploy/pkg/kube"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func main() {
	initLogger()
	log.Info().Msg("Hello World!")
	r := chi.NewRouter()

	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(requestIDMiddleware)
	r.Use(ChiLogger)
	r.Use(Recoverer)

	k, err := kube.New()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	r.Get("/api/v1/pods/list", func(w http.ResponseWriter, r *http.Request) {
		_, err := k.ListPods(r.Context())
		if err != nil {
			response.New().InternalServerError(w)
			return
		}
		response.New().Send(w, 200)
	})

	hserver := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	if err := hserver.ListenAndServe(); err != nil {
		log.Panic().Err(err).Send()
	}

}
