package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"kalos-cookbook/api"
	"kalos-cookbook/db"

	"github.com/lmittmann/tint"
	_ "modernc.org/sqlite"
)

func main() {
	ctx := context.Background()

	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelInfo,
			TimeFormat: time.RFC3339,
		}),
	))

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

	httpPort := "8080"
	slog.Info("HTTP server started", "port", httpPort)
	http.ListenAndServe(":"+httpPort, server.Router)
}
