package database

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/xid"
)

var ErrSessionNotFound = errors.New("session not found")

func (d *Database) CreateSession(ctx context.Context, userID uuid.UUID) (string, error) {
	sid := xid.New().String()
	_, err := d.P.Exec(ctx, `INSERT INTO users.sessions (session_id, user_id) VALUES ($1, $2)`, sid, userID)
	return sid, err
}

func (d *Database) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := d.P.Exec(ctx, `DELETE FROM users.sessions WHERE session_id = $1`, sessionID)
	return err
}

func (d *Database) SessionUserID(ctx context.Context, sessionID string) (uuid.UUID, error) {
	var id uuid.UUID
	err := d.P.QueryRow(ctx, `SELECT user_id FROM users.sessions WHERE session_id = $1`, sessionID).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return id, ErrSessionNotFound
	}
	return id, err
}