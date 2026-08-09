package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/auth"
	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/database"
	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/shared/check"
	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/shared/response"
	"github.com/matthewhartstonge/argon2"
	"github.com/rs/zerolog/log"
)

type RegisterRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

func Regsiter(d *database.Database, a *auth.Auth, argon argon2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := log.Ctx(ctx)

		req := new(RegisterRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			response.New().BadRequest(w)
			return
		}
		if err := check.Null(req); err != nil {
			response.New().Error(err.Error(), response.ErrorCodeNullField).Send(w, http.StatusBadRequest)
			return
		}

		password, err := argon.HashEncoded([]byte(req.Password))
		if err != nil {
			l.Err(err).Msg("password hash")
			response.New().InternalServerError(w)
			return
		}

		id, err := d.RegisterUser(ctx, database.RegisterUser{
			Email:         req.Email,
			UserFirstName: req.FirstName,
			UserLastName:  req.LastName,
			UserPassword:  string(password),
		})
		if err != nil {
			if errors.Is(err, database.ErrDuplicate) {
				response.New().Error("email alredy registered", "00000").Send(w, http.StatusConflict)
				return
			}
			l.Err(err).Msg("storing to database")
			response.New().InternalServerError(w)
			return
		}
		token, err := a.IssueSession(ctx, w, id)
		if err != nil {
			l.Err(err).Msg("issuing session")
			response.New().InternalServerError(w)
			return
		}
		if err := response.New().Response(map[string]string{
			"user_id":       id.String(),
			"session_token": token,
		}).
			Send(w, http.StatusCreated); err != nil {
			l.Err(err).Msg("writting response")
			return
		}
	}
}
