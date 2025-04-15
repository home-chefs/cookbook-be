package api

import (
	"context"
	"kalos-cookbook/db"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	DB     *db.DB
}

func NewServer(ctx context.Context, db *db.DB) *Server {
	server := &Server{
		DB: db,
	}

	r, err := NewRouter(ctx, *server)
	if err != nil {
		panic(err)
	}

	server.Router = r

	return server
}
