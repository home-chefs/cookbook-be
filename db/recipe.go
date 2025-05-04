package db

import (
	"context"
	"fmt"
	"log/slog"

	"kalos-cookbook/errors"
	"kalos-cookbook/types"
)

const (
	insertRecipeQuery    = `INSERT INTO recipe (name, time_to_cook, source, video_link, directions, cover_image_path) VALUES (?,?,?,?,?,?);`
	selectRecipeById     = `SELECT * FROM recipe WHERE id=?;`
	updateRecipeById     = `UPDATE recipe SET name=?, time_to_cook=?, source=?, video_link=?, directions=?, cover_image_path=? WHERE id=?;`
	deleteRecipeById     = `DELETE FROM recipe WHERE id=?;`
	selectRecipesBuilder = `SELECT * FROM recipe`
)

func (db *DB) CreateRecipe(ctx context.Context, recipe types.Recipe) (*types.Recipe, error) {
	// Insert recipe
	tx, err := db.sqlite.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "begin db tx")
	}
	res, err := tx.Exec(insertRecipeQuery, recipe.Name, recipe.TimeToCook, recipe.Source, recipe.VideoLink, recipe.Directions, recipe.CoverImagePath)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "exec db tx for recipe")
	}

	recipeID, err := res.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "fetch last inserted db id")
	}

	// Insert products
	for _, p := range recipe.Products {
		_, err = tx.Exec(insertProductQuery, p.Name, p.Amount, p.AmountType, recipeID)
		if err != nil {
			tx.Rollback()
			return nil, errors.Wrap(err, fmt.Sprintf("exec db tx for product %v", p.Name))
		}
	}

	// Insert labels
	for _, l := range recipe.Labels {
		_, err = tx.Exec(insertLabelQuery, l, recipeID)
		if err != nil {
			tx.Rollback()
			return nil, errors.Wrap(err, fmt.Sprintf("exec db tx for label %v", l))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrap(err, "commit db tx")
	}

	createdRecipe, err := db.GetRecipeByID(ctx, int(recipeID))
	if err != nil {
		return nil, errors.Wrap(err, "get freshly created recipe")
	}

	slog.Info("Created new recipe", "id", recipeID)
	return createdRecipe, nil
}

func (db *DB) GetRecipeByID(ctx context.Context, id int) (*types.Recipe, error) {
	var recipe types.Recipe
	row := db.sqlite.QueryRowContext(ctx, selectRecipeById, id)
	err := row.Scan(&recipe.ID, &recipe.Name, &recipe.TimeToCook, &recipe.Source, &recipe.VideoLink, &recipe.Directions, &recipe.CreatedAt, &recipe.CoverImagePath)
	if err != nil {
		return nil, errors.Wrap(err, "scan recipe row")
	}

	var products []types.Product
	rowsProducts, err := db.sqlite.QueryContext(ctx, selectProductsByRecipe, id)
	if err != nil {
		return nil, errors.Wrap(err, "query db for products by recipe")
	}
	defer rowsProducts.Close()

	for rowsProducts.Next() {
		var product types.Product

		if err := rowsProducts.Scan(
			&product.ID, &product.Name, &product.Amount, &product.AmountType, &product.ProductRecipe,
		); err != nil {
			return nil, errors.Wrap(err, "scan products rows")
		}
		products = append(products, product)
	}
	recipe.Products = products

	var labels []string
	rowsLabels, err := db.sqlite.QueryContext(ctx, selectLabelsByRecipe, id)
	if err != nil {
		return nil, errors.Wrap(err, "query db for labels by recipe")
	}
	defer rowsLabels.Close()

	for rowsLabels.Next() {
		var label types.Label

		if err := rowsLabels.Scan(
			&label.ID, &label.Label, &label.LabelRecipe,
		); err != nil {
			return nil, errors.Wrap(err, "scan labels rows")
		}
		labels = append(labels, label.Label)
	}
	recipe.Labels = labels

	return &recipe, nil
}

type GetRecipesOpts struct {
	SortBy string
}

func (db *DB) GetRecipes(ctx context.Context, opts *GetRecipesOpts) ([]types.Recipe, error) {
	query := selectRecipesBuilder

	if opts.SortBy != "" {
		query += " ORDER BY " + opts.SortBy
	}

	query += ";"

	var recipes []types.Recipe
	recipesRows, err := db.sqlite.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "query db for recipes")
	}

	for recipesRows.Next() {
		var recipe types.Recipe
		err = recipesRows.Scan(&recipe.ID, &recipe.Name, &recipe.TimeToCook, &recipe.Source, &recipe.VideoLink, &recipe.Directions, &recipe.CreatedAt, &recipe.CoverImagePath)
		if err != nil {
			return nil, errors.Wrap(err, "scan recipes rows")
		}
		recipes = append(recipes, recipe)
	}

	for idx, r := range recipes {
		products, err := db.GetProducts(ctx, GetProductOpts{RecipeID: r.ID})
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("get products for recipe with id %v", r.ID))
		}
		r.Products = products

		labels, err := db.GetLabels(ctx, GetLabelsOpts{RecipeID: r.ID})
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("get labels for recipe with id %v", r.ID))
		}
		r.Labels = labels

		recipes[idx] = r
	}

	return recipes, nil
}

func (db *DB) DeleteRecipeByID(ctx context.Context, id int) error {
	_, err := db.sqlite.ExecContext(ctx, deleteRecipeById, id)
	if err != nil {
		return errors.Wrap(err, "exec db tx for recipe")
	}

	return nil
}

func (db *DB) UpdateRecipe(ctx context.Context, recipe types.Recipe) error {
	tx, err := db.sqlite.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "begin db tx")
	}

	_, err = tx.Exec(updateRecipeById, recipe.Name, recipe.TimeToCook, recipe.Source, recipe.VideoLink, recipe.Directions, recipe.CoverImagePath, recipe.ID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "exec db tx for recipe")
	}

	_, err = tx.Exec(deleteProductByRecipe, recipe.ID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "exec db tx for product deletion")
	}
	for _, p := range recipe.Products {
		_, err = tx.Exec(
			insertProductQuery, p.Name, p.Amount, p.AmountType, recipe.ID,
		)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, fmt.Sprintf("exec db tx for product creation with recipe id %v", recipe.ID))
		}
	}

	_, err = tx.Exec(deleteLabelByRecipe, recipe.ID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "exec db tx for label deletion")
	}
	for _, l := range recipe.Labels {
		_, err = tx.Exec(insertLabelQuery, l, recipe.ID)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, fmt.Sprintf("exec db tx for label creation with recipe id %v", recipe.ID))
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "commit db tx")
	}

	slog.Info("Updated recipe", "id", recipe.ID)

	return nil
}
