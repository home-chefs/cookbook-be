package db

import (
	"context"

	"kalos-cookbook/errors"
	"kalos-cookbook/types"
)

const (
	insertLabelQuery     = `INSERT INTO label (name, labelrecipe) VALUES (?,?);`
	selectLabelsByRecipe = `SELECT * FROM label WHERE labelrecipe=?;`
	deleteLabelByRecipe  = `DELETE FROM label WHERE labelrecipe=?;`
	selectLabelsBuilder  = `SELECT * FROM label`
)

type GetLabelsOpts struct {
	RecipeID int
}

func (db *DB) GetLabels(ctx context.Context, opts GetLabelsOpts) ([]string, error) {
	query := selectLabelsBuilder
	var labels []string

	args := []any{}
	if opts.RecipeID != 0 {
		query += ` WHERE labelrecipe=?`
		args = append(args, opts.RecipeID)
	}

	query += ";"
	labelRows, err := db.sqlite.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query db for labels")
	}

	defer labelRows.Close()
	for labelRows.Next() {
		var label types.Label

		if err := labelRows.Scan(
			&label.ID, &label.Label, &label.LabelRecipe,
		); err != nil {
			return nil, errors.Wrap(err, "scan label row")
		}
		labels = append(labels, label.Label)
	}

	return labels, nil
}
