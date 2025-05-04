package api

import (
	"context"
	"log/slog"

	"kalos-cookbook/db"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	DB     *db.DB
}

func NewServer(ctx context.Context, db *db.DB) *Server {
	slog.Info("Starting new HTTP server...")
	server := &Server{
		DB: db,
	}

	server.Router = NewRouter(ctx, *server)

	return server
}
