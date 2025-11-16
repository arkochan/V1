package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresConnection(databaseURL string) (*pgxpool.Pool, error) {
	fmt.Print("HELOOOOO")
	fmt.Print(databaseURL)
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}
