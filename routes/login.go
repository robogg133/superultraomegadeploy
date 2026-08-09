package routes

import (
	"encoding/json"
	"net/http"

	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/auth"
	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/database"
	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/shared/check"
	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/shared/response"
	"github.com/matthewhartstonge/argon2"
	"github.com/rs/zerolog/log"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(d *database.Database, argon argon2.Config, a *auth.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := log.Ctx(ctx)

		req := new(LoginRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			response.New().BadRequest(w)
			return
		}
		if err := check.Null(req); err != nil {
			response.New().Error(err.Error(), response.ErrorCodeNullField).Send(w, http.StatusBadRequest)
			return
		}

		u, err := d.GetUserByEmail(ctx, req.Email)
		if err != nil {
			response.New().Error("invalid credentials", response.ErrorCodeUnauthorized).Send(w, http.StatusUnauthorized)
			return
		}
		ok, err := argon2.VerifyEncoded([]byte(req.Password), []byte(u.Password))
		if err != nil || !ok {
			response.New().Error("invalid credentials", response.ErrorCodeUnauthorized).Send(w, http.StatusUnauthorized)
			return
		}

		token, err := a.IssueSession(ctx, w, u.ID)
		if err != nil {
			l.Err(err).Msg("issuing session")
			response.New().InternalServerError(w)
			return
		}
		response.New().Response(map[string]string{
			"user_id":      u.ID.String(),
			"session_token": token,
		}).Send(w, http.StatusOK)
	}
}