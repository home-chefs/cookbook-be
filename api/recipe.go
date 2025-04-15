package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kalos-cookbook/db"
	"kalos-cookbook/types"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func (s *Server) createRecipe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, err)
		return
	}

	var recipe types.Recipe
	err = json.Unmarshal(body, &recipe)
	if err != nil {
		badRequest(w, err)
		return
	}

	recipeCreated, err := s.DB.CreateRecipe(ctx, recipe)
	if err != nil {
		badRequest(w, err)
		return
	}

	recipeJSON, err := json.Marshal(recipeCreated)
	if err != nil {
		badRequest(w, err)
		return
	}

	w.Write(recipeJSON)
}

func (s *Server) getAllRecipes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sortBy := r.URL.Query().Get("sort_by")
	var err error

	if sortBy != "" {
		sortBy, err = formatSort(sortBy)
		if err != nil {
			badRequest(w, err)
			return
		}
	}

	labelFilter := r.URL.Query().Get("label")
	productFilter := r.URL.Query().Get("product")

	opts := &db.GetRecipesOpts{
		SortBy: sortBy,
	}
	recipes, err := s.DB.GetRecipes(ctx, opts)
	if err != nil {
		badRequest(w, err)
		return
	}

	filteredRecipes := []types.Recipe{}
	for _, r := range recipes {
		if recipeContainsProduct(r, productFilter) && recipeContainsLabel(r, labelFilter) {
			filteredRecipes = append(filteredRecipes, r)
		}
	}

	recipesJSON, err := json.Marshal(filteredRecipes)
	if err != nil {
		badRequest(w, err)
		return
	}
	w.Write(recipesJSON)
}

func (s *Server) getRecipe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	idVar := vars["id"]
	id, err := strconv.Atoi(idVar)
	if err != nil {
		badRequest(w, err)
		return
	}

	recipe, err := s.DB.GetRecipeByID(ctx, id)
	if err != nil {
		badRequest(w, err)
		return
	}
	recipeJSON, err := json.Marshal(recipe)
	if err != nil {
		badRequest(w, err)
		return
	}
	w.Write(recipeJSON)
}

func (s *Server) deleteRecipe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	idVar := vars["id"]
	id, err := strconv.Atoi(idVar)
	if err != nil {
		badRequest(w, err)
		return
	}

	err = s.DB.DeleteRecipeByID(ctx, id)
	if err != nil {
		badRequest(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) updateRecipe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, err)
		return
	}

	var recipe types.Recipe
	err = json.Unmarshal(body, &recipe)
	if err != nil {
		badRequest(w, err)
		return
	}

	err = s.DB.UpdateRecipe(ctx, recipe)
	if err != nil {
		badRequest(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func formatSort(sortBy string) (string, error) {
	splits := strings.Split(sortBy, ".")
	if len(splits) != 2 {
		return "", errors.New("malformed sortBy query parameter, should be field.order_direction")
	}
	field, order := splits[0], splits[1]
	if order != "desc" && order != "asc" {
		return "", errors.New("malformed orderdirection in sortBy query parameter, should be asc or desc")
	}
	if !slices.Contains(filterableFields(), field) {
		return "", errors.New("unknown field in sortBy query parameter")
	}
	return fmt.Sprintf("%s %s", field, strings.ToUpper(order)), nil
}

func filterableFields() []string {
	return []string{"created_at", "name", "time_to_cook"}
}

func recipeContainsProduct(recipe types.Recipe, product string) bool {
	if product == "" {
		return true
	}

	for _, p := range recipe.Products {
		if p.Name == product {
			return true
		}
	}
	return false
}

func recipeContainsLabel(recipe types.Recipe, label string) bool {
	if label == "" {
		return true
	}

	return slices.Contains(recipe.Labels, label)
}
