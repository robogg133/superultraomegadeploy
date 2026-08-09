package database

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type RegisterUser struct {
	Email         string
	UserFirstName string
	UserLastName  string
	UserPassword  string
}

func (d *Database) RegisterUser(ctx context.Context, u RegisterUser) (userID uuid.UUID, err error) {
	err = d.P.QueryRow(ctx, `INSERT INTO users.user (
		user_email,
		user_first_name,
		user_last_name,
		user_password
		)
		VALUES ($1, $2, $3, $4)
		RETURNING user_id`,
		u.Email,
		u.UserFirstName,
		u.UserLastName,
		u.UserPassword,
	).Scan(&userID)
	if pgerr, ok := errors.AsType[*pgconn.PgError](err); ok {
		if pgerr.Code == "2305" {
			return uuid.Nil, ErrDuplicate
		}
	}
	return userID, err
}
