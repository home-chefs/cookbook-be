package main

import (
	"context"
	"kalos-cookbook/api"
	"kalos-cookbook/db"
	"net/http"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func main() {
	ctx := context.Background()

	err := os.MkdirAll("db", 0o777)
	if err != nil && err != os.ErrExist {
		panic(err)
	}
	path := filepath.Join(".", "sqlite", "sqlite.db")
	sqlite, err := db.InitSqlite(ctx, path)
	if err != nil {
		panic(err)
	}

	err = sqlite.RunMigrations(ctx)
	if err != nil {
		panic(err)
	}

	server := api.NewServer(ctx, sqlite)

	http.ListenAndServe(":8080", server.Router)
}
