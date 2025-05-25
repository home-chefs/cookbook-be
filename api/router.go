package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(ctx context.Context, server Server) *mux.Router {
	r := mux.NewRouter()

	r.Use(contentTypeApplicationJsonMiddleware)
	r.Use(accessControlMiddleware)

	r.HandleFunc("/health", server.health).Methods(http.MethodGet)

	r.HandleFunc("/recipes", server.createRecipe).Methods(http.MethodPost)
	r.HandleFunc("/recipes", server.getAllRecipes).Methods(http.MethodGet)
	r.HandleFunc("/recipes/{id}", server.getRecipe).Methods(http.MethodGet)
	r.HandleFunc("/recipes/{id}", server.deleteRecipe).Methods(http.MethodDelete)
	r.HandleFunc("/recipes", server.updateRecipe).Methods(http.MethodPut)

	return r
}

func contentTypeApplicationJsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func accessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Credentials", "true")
		header.Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		header.Set("Access-Control-Allow-Headers", "*")
		header.Set("Access-Control-Max-Age", "86400")

		next.ServeHTTP(w, r)
	})
}
