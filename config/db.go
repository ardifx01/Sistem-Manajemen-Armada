package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func InitPostgres() (*pgx.Conn, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	var conn *pgx.Conn
	var err error

	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		conn, err = pgx.Connect(ctx, dsn)
		cancel()

		if err == nil {
			log.Println(" Connected to Postgres (pgx)")

			// create table if not exists
			_, err = conn.Exec(context.Background(), `
				CREATE TABLE IF NOT EXISTS vehicle_locations (
					vehicle_id TEXT NOT NULL,
					latitude DOUBLE PRECISION NOT NULL,
					longitude DOUBLE PRECISION NOT NULL,
					timestamp TIMESTAMP NOT NULL
				)
			`)
			if err != nil {
				return nil, fmt.Errorf("failed to create table: %w", err)
			}
			log.Println("Table vehicle_locations ready")
			return conn, nil
		}

		log.Printf("Waiting for Postgres... (%d/10): %v", i+1, err)
		time.Sleep(3 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect Postgres after retries: %w", err)
}
