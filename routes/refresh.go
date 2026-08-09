package routes

import (
	"net/http"

	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/auth"
	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/shared/response"
	"github.com/rs/zerolog/log"
)

// @Summary Refresh the short-lived session token
// @Tags auth
// @Produce json
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/refresh [post]
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