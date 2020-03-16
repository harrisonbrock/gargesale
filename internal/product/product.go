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

func Retrieve(db *sqlx.DB, id string) (*Product, error) {

	// Create a slice of products.
	var p Product

	const q = `SELECT product_id, name, cost, quantity, date_updated, date_created 
	FROM products 
	WHERE product_id = $1`

	if err := db.Get(&p, q, id); err != nil {
		return nil, err
	}

	return &p, nil
}
