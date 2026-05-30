package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(databaseURL)

	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v", err)
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to create database pool: %v", err)
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		pool.Close()
		return nil, err
	}

	log.Println("Successfully connected to the database")

	return pool, nil
}
