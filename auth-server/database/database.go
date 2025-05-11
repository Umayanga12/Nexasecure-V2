// internal/database/database.go
package database

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect() error {

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	connStr := fmt.Sprintf(
        "postgresql://%s:%s@%s:%s/%s",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return err
	}

	DB = pool
	return nil
}

func CreateTables() error {
	_, err := DB.Exec(context.Background(), `
		CREATE EXTENSION IF NOT EXISTS "pgcrypto";

		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			reqpubaddr VARCHAR(300) NOT NULL,
			authpubaddr VARCHAR(300) NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS otps (
			user_id UUID REFERENCES users(id),
			code VARCHAR(6) NOT NULL,
			expires_at TIMESTAMP NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("error creating tables: %w", err)
	}
	return nil
}