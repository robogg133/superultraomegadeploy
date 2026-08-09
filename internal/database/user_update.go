package database

import (
	"context"

	"github.com/google/uuid"
)

func (d *Database) UpdateUser(ctx context.Context, userID uuid.UUID, u User) error {
	_, err := d.P.Exec(ctx, `UPDATE users.user SET
		user_email = $2,
		user_first_name = $3,
		user_last_name = $4,
		user_password = $5
		WHERE user_id = $1`,
		userID,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Password,
	)
	return err
}