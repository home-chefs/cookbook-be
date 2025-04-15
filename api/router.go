package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(ctx context.Context, server Server) (*mux.Router, error) {
	r := mux.NewRouter()

	r.Use(contentTypeApplicationJsonMiddleware)
	r.Use(accessControlMiddleware)

	r.HandleFunc("/recipes", server.createRecipe).Methods(http.MethodPost)
	r.HandleFunc("/recipes", server.getAllRecipes).Methods(http.MethodGet)
	r.HandleFunc("/recipes/{id}", server.getRecipe).Methods(http.MethodGet)
	r.HandleFunc("/recipes/{id}", server.deleteRecipe).Methods(http.MethodDelete)
	r.HandleFunc("/recipes", server.updateRecipe).Methods(http.MethodPut)

	return r, nil
}

func contentTypeApplicationJsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func accessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
