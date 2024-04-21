package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

type PostgresDB struct {
	*sql.DB
}

func NewPostgresDB(dbURL string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	if err := migrateDB(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	return &PostgresDB{db}, nil
}

func migrateDB(db *sql.DB) error {
	goose.SetDialect("postgres")

	if err := goose.Run("up", db, "pkg/db/migration"); err != nil {
		return fmt.Errorf("failed to run database migrations: %v", err)
	}

	return nil
}
