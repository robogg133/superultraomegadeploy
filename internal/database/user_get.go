package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	FirstName string
	LastName  string
	Password  string
	CreatedAt time.Time
}

func (d *Database) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := new(User)
	err := d.P.QueryRow(ctx, `SELECT
			user_id,
			user_email,
			user_first_name,
			user_last_name,
			user_password,
			created_at
			FROM users.user
			WHERE user_email = $1`, email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.FirstName,
		&u.LastName,
		&u.Password,
		&u.CreatedAt,
	)
	return u, err
}

func (d *Database) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	u := new(User)
	err := d.P.QueryRow(ctx, `SELECT
		user_id,
		user_email,
		user_first_name,
		user_last_name,
		user_password,
		created_at
		FROM users.user
		WHERE user_id = $1`, userID,
	).Scan(
		&u.ID,
		&u.Email,
		&u.FirstName,
		&u.LastName,
		&u.Password,
		&u.CreatedAt,
	)
	return u, err
}
