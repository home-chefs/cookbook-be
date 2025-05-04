package db

import (
	"context"
	"database/sql"
	"kalos-cookbook/errors"
	"log/slog"
)

type DB struct {
	sqlite *sql.DB
}

func InitSqlite(ctx context.Context, dbPath string) (*DB, error) {
	slog.Info("Initialiasing SQLite database...", "path", dbPath)
	var err error
	sqlite, err := sql.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, errors.Wrap(err, "open SQL database")
	}
	_, err = sqlite.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS db_metadata (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			migration TEXT NOT NULL UNIQUE
			);`,
	)
	if err != nil {
		return nil, errors.Wrap(err, "create db_metadata table")
	}
	dbWrap := DB{sqlite: sqlite}

	slog.Info("Initialised SQLite database")
	return &dbWrap, nil
}
