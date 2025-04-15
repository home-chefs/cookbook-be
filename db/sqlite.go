package db

import (
	"context"
	"database/sql"
)

type DB struct {
	sqlite *sql.DB
}

func InitSqlite(ctx context.Context, dbPath string) (*DB, error) {
	var err error
	sqlite, err := sql.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, err
	}
	_, err = sqlite.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS db_metadata (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			migration TEXT NOT NULL UNIQUE
			);`,
	)
	if err != nil {
		return nil, err
	}
	dbWrap := DB{sqlite: sqlite}

	return &dbWrap, nil
}
