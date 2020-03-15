package product

import (
	"github.com/jmoiron/sqlx"
)

// List returns all know Products.
func List(db *sqlx.DB) ([]Product, error) {

	// Create a slice of products.
	list := []Product{}

	const q = `SELECT * FROM products`

	if err := db.Select(&list, q); err != nil {
		return nil, err
	}

	return list, nil
}
