package types

type Product struct {
	ID         int    `json:"-"`
	Name       string `json:"name"`
	Amount     int    `json:"amount"`
	AmountType string `json:"amount_type"`

	ProductRecipe int `json:"-"`
}
