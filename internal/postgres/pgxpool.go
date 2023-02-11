package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPgxPoolConnection() (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(context.Background(), "postgres://root:secret@localhost/gosh")
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable ping to database: %v", err)
	}

	return pool, nil
}
