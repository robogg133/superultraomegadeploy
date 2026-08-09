package database

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrForbidden = errors.New("permission denied")

const sqlStateInsufficientPrivilege = "42501"

// SetServerConfig stores a server config key only if the user has the
// can_change_server_settings permission, enforced by the database itself.
func (d *Database) SetServerConfig(ctx context.Context, userID uuid.UUID, key string, value json.RawMessage) error {
	_, err := d.P.Exec(ctx, `CALL set_server_config($1, $2, $3)`, userID, key, value)
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == sqlStateInsufficientPrivilege {
		return ErrForbidden
	}
	return err
}