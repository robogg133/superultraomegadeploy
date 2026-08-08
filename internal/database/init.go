package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	P *pgxpool.Pool
}

func Init(ctx context.Context, connString string) (*Database, error) {
	d := new(Database)
	var err error
	d.P, err = pgxpool.New(ctx, connString)
	return d, err
}
