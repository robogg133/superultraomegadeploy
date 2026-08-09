package main

import (
	"context"
	"embed"
	"io"
	"net/http"

	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/auth"
	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/database"
	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/shared/response"
	"git.servidordomal.lol/robogg133/superultraomegadeploy/pkg/kube"
	"git.servidordomal.lol/robogg133/superultraomegadeploy/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/matthewhartstonge/argon2"
	"github.com/rs/zerolog/log"
)

//go:embed include/**
var includeFs embed.FS

func main() {
	appCtx, appCtxCancel := context.WithCancel(context.Background())
	defer appCtxCancel()
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
	databaseCtx := context.WithoutCancel(appCtx)

	d, err := database.Init(databaseCtx, PostgresConnString)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	argon := argon2.RecommendedDefaults()
	a := auth.New(d, JWTSecret)

	r.Post("/api/v1/users/register", routes.Regsiter(d, a, argon))
	r.Post("/api/v1/auth/login", routes.Login(d, argon, a))
	r.Post("/api/v1/auth/refresh", routes.Refresh(a))

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

	r.With(a.RequireAuth).Get("/api/v1/pods/list", routes.PodsList(k))

	hserver := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	if err := hserver.ListenAndServe(); err != nil {
		log.Panic().Err(err).Send()
	}

}
