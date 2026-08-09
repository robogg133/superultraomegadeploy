package database

import (
	"context"

	"github.com/google/uuid"

	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/shared/userperm"
)

func (d *Database) UserCount(ctx context.Context) (int64, error) {
	var count int64
	err := d.P.QueryRow(ctx, `SELECT COUNT(*) FROM users.user`).Scan(&count)
	return count, err
}

func (d *Database) GrantAllPermissions(ctx context.Context, userID uuid.UUID) error {
	for _, perm := range userperm.All {
		if _, err := d.P.Exec(ctx, `INSERT INTO users.permissions (user_id, permission)
			VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, perm); err != nil {
			return err
		}
	}
	return nil
}