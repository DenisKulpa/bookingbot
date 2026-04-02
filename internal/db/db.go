package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func New(dataSourceName string) (*sql.DB, error) {
	database, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("db.New: %%w", err)
	}

	database.SetMaxOpenConns(1)
	database.SetMaxIdleConns(1)
	database.SetConnMaxLifetime(0)

	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %%w", err)
	}

	if _, err := database.Exec(`PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;`); err != nil {
		return nil, fmt.Errorf("db pragma: %%w", err)
	}

	migration, err := os.ReadFile("migrations/000001_create_zones.up.sql")
	if err != nil {
		return nil, fmt.Errorf("db: read migration: %%w", err)
	}
	if _, err := database.Exec(string(migration)); err != nil {
		return nil, fmt.Errorf("db: apply migration: %%w", err)
	}

	return database, nil
}
