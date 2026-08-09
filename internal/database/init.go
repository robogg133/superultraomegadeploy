package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	P *pgxpool.Pool
}

func Init(ctx context.Context, connString string) (*Database, error) {
	d := new(Database)
	var err error
	d.P, err = pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	pCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	return d, d.P.Ping(pCtx)
}
