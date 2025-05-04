package db

import (
	"context"
	"fmt"
	"kalos-cookbook/errors"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
)

func (db *DB) RunMigrations(ctx context.Context) error {
	migrationsPath := "./db/migrate/*.sql"
	slog.Info("Running SQLite migrations...", "path", migrationsPath)
	migrationsDB, err := db.readMetadata(ctx)
	if err != nil {
		return errors.Wrap(err, "read sqlite metadata")
	}

	migrationsFiles, err := filepath.Glob(migrationsPath)
	if err != nil {
		return errors.Wrap(err, "fetch migration files from path")
	}

	for _, file := range migrationsFiles {
		filename := filepath.Base(file)

		if !slices.Contains(migrationsDB, filename) {
			slog.Info("Running new migration", "path", filename)
			query, err := os.ReadFile(file)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("fetch migration files from path for migration %v", filename))
			}
			tx, err := db.sqlite.BeginTx(ctx, nil)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("begin db tx for migration %v", filename))
			}

			_, err = tx.Exec(string(query))
			if err != nil {
				tx.Rollback()
				return errors.Wrap(err, fmt.Sprintf("exec db migration for migration %v", filename))
			}

			_, err = tx.Exec(`INSERT INTO db_metadata (migration) VALUES (?);`, filename)
			if err != nil {
				tx.Rollback()
				return errors.Wrap(err, fmt.Sprintf("exec migration in db_metadata for migration %v", filename))
			}

			err = tx.Commit()
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("commit db tx for migration %v", filename))
			}
		}
	}

	slog.Info("SQLite migration done")
	return nil
}

func (db *DB) readMetadata(ctx context.Context) ([]string, error) {
	var migrations []string

	rows, err := db.sqlite.QueryContext(ctx, `SELECT migration FROM db_metadata;`)
	if err != nil {
		return nil, errors.Wrap(err, "query migration from db_metadata")
	}

	defer rows.Close()
	for rows.Next() {
		var migration string
		err := rows.Scan(&migration)
		if err != nil {
			return nil, errors.Wrap(err, "scan migration for row")
		}
		migrations = append(migrations, migration)
	}

	return migrations, nil
}
