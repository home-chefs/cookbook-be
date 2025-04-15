package db

import (
	"context"
	"kalos-cookbook/types"
)

const (
	insertProductQuery     = `INSERT INTO product (name, amount, amount_type, productrecipe) VALUES (?,?,?,?);`
	selectProductsByRecipe = `SELECT * FROM product WHERE productrecipe=?;`
	deleteProductByRecipe  = `DELETE FROM product WHERE productrecipe=?;`
	selectProductsBuilder  = `SELECT * FROM product`
)

type GetProductOpts struct {
	RecipeID int
}

func (db *DB) GetProducts(ctx context.Context, opts GetProductOpts) ([]types.Product, error) {
	query := selectProductsBuilder
	var products []types.Product

	args := []any{}
	if opts.RecipeID != 0 {
		query += ` WHERE productrecipe=?`
		args = append(args, opts.RecipeID)
	}

	query += ";"
	productRows, err := db.sqlite.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer productRows.Close()
	for productRows.Next() {
		var product types.Product

		if err := productRows.Scan(
			&product.ID, &product.Name, &product.Amount, &product.AmountType, &product.ProductRecipe,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
