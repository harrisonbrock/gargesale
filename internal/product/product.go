package product

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

var (
	ErrNotFound  = errors.New("product not found")
	ErrInvalidId = errors.New("id provided not a valid UUID")
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

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidId
	}

	// Create a slice of products.
	var p Product

	const q = `SELECT product_id, name, cost, quantity, date_updated, date_created 
	FROM products 
	WHERE product_id = $1`

	if err := db.Get(&p, q, id); err != nil {

		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

// Create makes a new Product.
func Create(db *sqlx.DB, np NewProduct, now time.Time) (*Product, error) {
	p := Product{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		DateCreated: now,
		DateUpdated: now,
	}
	const q = `INSERT INTO products
	(product_id, name, cost, quantity, date_created, date_updated)
	VALUES($1, $2, $3, $4, $5, $6)`

	if _, err := db.Exec(q, p.ID, p.Name, p.Cost, p.Quantity, p.DateCreated, p.DateUpdated); err != nil {
		return nil, errors.Wrap(err, "inserting product")
	}
	return &p, nil
}
