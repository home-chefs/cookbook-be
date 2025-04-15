package db

import (
	"context"
	"os"
	"path/filepath"
	"slices"
)

func (db *DB) RunMigrations(ctx context.Context) error {
	migrationsDB, err := db.readMetadata(ctx)
	if err != nil {
		return err
	}

	migrationsFiles, err := filepath.Glob("./db/migrate/*.sql")
	if err != nil {
		return err
	}

	for _, file := range migrationsFiles {
		filename := filepath.Base(file)

		if !slices.Contains(migrationsDB, filename) {
			query, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			tx, err := db.sqlite.BeginTx(ctx, nil)
			if err != nil {
				return err
			}

			_, err = tx.Exec(string(query))
			if err != nil {
				tx.Rollback()
				return err
			}

			_, err = tx.Exec(`INSERT INTO db_metadata (migration) VALUES (?);`, filename)
			if err != nil {
				tx.Rollback()
				return err
			}

			err = tx.Commit()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (db *DB) readMetadata(ctx context.Context) ([]string, error) {
	var migrations []string

	rows, err := db.sqlite.QueryContext(ctx, `SELECT migration FROM db_metadata;`)
	if err != nil {
		return migrations, err
	}

	defer rows.Close()
	for rows.Next() {
		var migration string
		err := rows.Scan(&migration)
		if err != nil {
			return migrations, err
		}
		migrations = append(migrations, migration)
	}

	return migrations, nil
}
