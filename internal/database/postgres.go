package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	var ctx context.Context = context.Background()

	var cfg *pgxpool.Config
	var err error
	cfg, err = pgxpool.ParseConfig(databaseURL)

	if err != nil {
		log.Printf("Unable to parse DATABASE_URL: %v", err)
	}

	var pool *pgxpool.Pool
	pool, err = pgxpool.NewWithConfig(ctx, cfg)

	if err != nil {
		log.Printf("Unable to create connection pool: %v", err)
		return nil, err
	}

	err = pool.Ping(ctx)

	if err != nil {
		log.Printf("Unable to ping database: %v", err)
		pool.Close()
		return nil, err
	}

	log.Printf("Successfully connected to PostgreSQL database")

	return pool, nil
}
