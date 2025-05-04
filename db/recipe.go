package db

import (
	"context"
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
		return nil, err
	}
	res, err := tx.Exec(insertRecipeQuery, recipe.Name, recipe.TimeToCook, recipe.Source, recipe.VideoLink, recipe.Directions, recipe.CoverImagePath)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	recipeID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Insert products
	for _, p := range recipe.Products {
		_, err = tx.Exec(insertProductQuery, p.Name, p.Amount, p.AmountType, recipeID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Insert labels
	for _, l := range recipe.Labels {
		_, err = tx.Exec(insertLabelQuery, l, recipeID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	createdRecipe, err := db.GetRecipeByID(ctx, int(recipeID))
	if err != nil {
		return nil, err
	}

	return createdRecipe, nil
}

func (db *DB) GetRecipeByID(ctx context.Context, id int) (*types.Recipe, error) {
	var recipe types.Recipe
	row := db.sqlite.QueryRowContext(ctx, selectRecipeById, id)
	err := row.Scan(&recipe.ID, &recipe.Name, &recipe.TimeToCook, &recipe.Source, &recipe.VideoLink, &recipe.Directions, &recipe.CreatedAt, &recipe.CoverImagePath)
	if err != nil {
		return nil, err
	}

	var products []types.Product
	rowsProducts, err := db.sqlite.QueryContext(ctx, selectProductsByRecipe, id)
	if err != nil {
		return nil, err
	}
	defer rowsProducts.Close()

	for rowsProducts.Next() {
		var product types.Product

		if err := rowsProducts.Scan(
			&product.ID, &product.Name, &product.Amount, &product.AmountType, &product.ProductRecipe,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	recipe.Products = products

	var labels []string
	rowsLabels, err := db.sqlite.QueryContext(ctx, selectLabelsByRecipe, id)
	if err != nil {
		return nil, err
	}
	defer rowsLabels.Close()

	for rowsLabels.Next() {
		var label types.Label

		if err := rowsLabels.Scan(
			&label.ID, &label.Label, &label.LabelRecipe,
		); err != nil {
			return nil, err
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
		return nil, err
	}

	for recipesRows.Next() {
		var recipe types.Recipe
		err = recipesRows.Scan(&recipe.ID, &recipe.Name, &recipe.TimeToCook, &recipe.Source, &recipe.VideoLink, &recipe.Directions, &recipe.CreatedAt, &recipe.CoverImagePath)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}

	for idx, r := range recipes {
		products, err := db.GetProducts(ctx, GetProductOpts{RecipeID: r.ID})
		if err != nil {
			return nil, err
		}
		r.Products = products

		labels, err := db.GetLabels(ctx, GetLabelsOpts{RecipeID: r.ID})
		if err != nil {
			return nil, err
		}
		r.Labels = labels

		recipes[idx] = r
	}

	return recipes, nil
}

func (db *DB) DeleteRecipeByID(ctx context.Context, id int) error {
	_, err := db.sqlite.ExecContext(ctx, deleteRecipeById, id)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateRecipe(ctx context.Context, recipe types.Recipe) error {
	tx, err := db.sqlite.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.Exec(updateRecipeById, recipe.Name, recipe.TimeToCook, recipe.Source, recipe.VideoLink, recipe.Directions, recipe.CoverImagePath, recipe.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(deleteProductByRecipe, recipe.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, p := range recipe.Products {
		_, err = tx.Exec(
			insertProductQuery, p.Name, p.Amount, p.AmountType, recipe.ID,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	_, err = tx.Exec(deleteLabelByRecipe, recipe.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, l := range recipe.Labels {
		_, err = tx.Exec(insertLabelQuery, l, recipe.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
