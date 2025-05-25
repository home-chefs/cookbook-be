package api

import (
	"context"
	"log/slog"
	"net/http"

	"kalos-cookbook/db"

	"github.com/rs/cors"
)

type Server struct {
	Handler http.Handler
	DB      *db.DB
}

func NewServer(ctx context.Context, db *db.DB) *Server {
	slog.Info("Starting new HTTP server...")
	server := &Server{
		DB: db,
	}

	router := NewRouter(ctx, *server)
	handler := cors.AllowAll().Handler(router)
	server.Handler = handler

	return server
}
