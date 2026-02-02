package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Parse connection config for pgx
	config, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection config: %w", err)
	}

	// Create stdlib connector
	stdlib.RegisterConnConfig(config)
	db := stdlib.OpenDB(*config)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	// NOTE: DB schema (tables and relations) is managed externally/manual.
	// This function only opens a connection and verifies it. No automatic
	// migrations or ALTER TABLE statements are performed here.

	log.Println("Database connected successfully")
	return db, nil
}
