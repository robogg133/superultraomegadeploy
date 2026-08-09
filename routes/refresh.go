package routes

import (
	"net/http"

	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/auth"
	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/shared/response"
	"github.com/rs/zerolog/log"
)

func Refresh(a *auth.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := log.Ctx(ctx)

		c, err := r.Cookie(auth.PersistCookie)
		if err != nil {
			response.New().Error("unauthorized", response.ErrorCodeUnauthorized).Send(w, http.StatusUnauthorized)
			return
		}

		token, err := a.RefreshSession(ctx, w, c.Value)
		if err != nil {
			l.Err(err).Msg("refreshing session")
			response.New().Error("unauthorized", response.ErrorCodeUnauthorized).Send(w, http.StatusUnauthorized)
			return
		}
		response.New().Response(map[string]string{"session_token": token}).Send(w, http.StatusOK)
	}
}