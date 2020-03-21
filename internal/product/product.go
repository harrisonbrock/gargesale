package product

import (
	"context"
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
func List(ctx context.Context, db *sqlx.DB) ([]Product, error) {

	// Create a slice of products.
	list := []Product{}

	const q = `SELECT * FROM products`

	if err := db.SelectContext(ctx, &list, q); err != nil {
		return nil, err
	}

	return list, nil
}

func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*Product, error) {

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidId
	}

	// Create a slice of products.
	var p Product

	const q = `SELECT product_id, name, cost, quantity, date_updated, date_created 
	FROM products 
	WHERE product_id = $1`

	if err := db.GetContext(ctx, &p, q, id); err != nil {

		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

// Create makes a new Product.
func Create(ctx context.Context, db *sqlx.DB, np NewProduct, now time.Time) (*Product, error) {
	p := Product{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}
	const q = `INSERT INTO products
	(product_id, name, cost, quantity, date_created, date_updated)
	VALUES($1, $2, $3, $4, $5, $6)`

	if _, err := db.ExecContext(ctx, q, p.ID, p.Name, p.Cost, p.Quantity, p.DateCreated, p.DateUpdated); err != nil {
		return nil, errors.Wrap(err, "inserting product")
	}
	return &p, nil
}
