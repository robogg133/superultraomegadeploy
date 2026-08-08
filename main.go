package main

import (
	"embed"
	"io"
	"net/http"

	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/shared/response"
	"git.servidordomal.lol/robogg133/superultraomegadeploy/pkg/kube"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

//go:embed include/**
var includeFs embed.FS

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
	r.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		f, err := includeFs.Open("include/robots.txt")
		if err != nil {
			log.Ctx(r.Context()).Err(err).Send()
			response.New().InternalServerError(w)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		if _, err := io.Copy(w, f); err != nil {
			log.Ctx(r.Context()).Err(err).Send()
			response.New().InternalServerError(w)
			return
		}
	})

	r.Get("/api/v1/pods/list", func(w http.ResponseWriter, r *http.Request) {
		list, err := k.ListPods(r.Context())
		if err != nil {
			log.Ctx(r.Context()).Err(err).Send()
			response.New().InternalServerError(w)
			return
		}
		resp := make([]map[string]any, 0)

		for _, pod := range list {
			r := make(map[string]any)
			r["name"] = pod.Name
			r["namespace"] = pod.Namespace
			r["creation_time"] = pod.CreationTimestamp.Time
			resp = append(resp, r)
		}
		log.Ctx(r.Context()).Debug().Int("len_pods", len(resp)).Send()

		response.New().Response(resp).Send(w, 200)
	})

	hserver := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	if err := hserver.ListenAndServe(); err != nil {
		log.Panic().Err(err).Send()
	}

}
