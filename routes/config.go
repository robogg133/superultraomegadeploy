package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/auth"
	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/database"
	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/shared/check"
	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/shared/response"
	"github.com/rs/zerolog/log"
)

type SetConfigRequest struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

// @Summary Set a server config
// @Tags config
// @Accept json
// @Produce json
// @Param body body SetConfigRequest true "config key and value"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/config [post]
func SetConfig(d *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := log.Ctx(ctx)

		userID, ok := auth.UserID(ctx)
		if !ok {
			response.New().Error("unauthorized", response.ErrorCodeUnauthorized).Send(w, http.StatusUnauthorized)
			return
		}

		req := new(SetConfigRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			response.New().BadRequest(w)
			return
		}
		if err := check.Null(req); err != nil {
			response.New().Error(err.Error(), response.ErrorCodeNullField).Send(w, http.StatusBadRequest)
			return
		}

		value, err := json.Marshal(req.Value)
		if err != nil {
			response.New().BadRequest(w)
			return
		}

		if err := d.SetServerConfig(ctx, userID, req.Key, value); err != nil {
			if errors.Is(err, database.ErrForbidden) {
				response.New().Error("forbidden", response.ErrorCodeForbidden).Send(w, http.StatusForbidden)
				return
			}
			l.Err(err).Msg("setting server config")
			response.New().InternalServerError(w)
			return
		}
		response.New().Response(map[string]bool{"ok": true}).Send(w, http.StatusOK)
	}
}