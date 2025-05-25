package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(ctx context.Context, server Server) *mux.Router {
	r := mux.NewRouter()

	r.Use(contentTypeApplicationJsonMiddleware)

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
