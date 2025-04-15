package types

type Label struct {
	ID    int    `json:"-"`
	Label string `json:"label"`

	LabelRecipe int `json:"-"`
}
