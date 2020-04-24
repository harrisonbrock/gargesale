package product

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/harrisonbrock/gargesale/internal/platform/auth"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"time"
)

var (
	ErrNotFound  = errors.New("product not found")
	ErrInvalidId = errors.New("id provided not a valid UUID")
	// ErrForbidden occurs when a user tries to do something that is forbidden to
	// them according to our access control policies.
	ErrForbidden = errors.New("Attempted action is not allowed")
)

// List returns all know Products.
func List(ctx context.Context, db *sqlx.DB) ([]Product, error) {

	// Create a slice of products.
	list := []Product{}

	const q = `SELECT 
	p.product_id, p.name, p.cost, p.quantity,
    coalesce(sum(s.quantity),0) as sold, coalesce(sum(s.paid), 0) as revenue,
    p.date_created, p.date_updated
	FROM products as p
	LEFT JOIN sales As s on p.product_id = s.product_id
	group by p.product_id
	`

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

	const q = `SELECT 
       p.product_id, p.name, p.cost, p.quantity,
       coalesce(sum(s.quantity), 0) as sold, coalesce(sum(s.paid), 0) as revenue,
       p.date_updated, p.date_created
	FROM products as p 
	left join sales s on p.product_id = s.product_id
	WHERE p.product_id = $1
	group by p.product_id
	`

	if err := db.GetContext(ctx, &p, q, id); err != nil {

		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

// Create makes a new Product.
func Create(ctx context.Context, db *sqlx.DB, user auth.Claims, np NewProduct, now time.Time) (*Product, error) {
	p := Product{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Cost:        np.Cost,
		UserID:      user.Subject,
		Quantity:    np.Quantity,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}
	const q = `INSERT INTO products
	(product_id, name, cost, quantity, user_id, date_created, date_updated)
	VALUES($1, $2, $3, $4, $5, $6, $7)`

	if _, err := db.ExecContext(ctx, q, p.ID, p.Name, p.Cost, p.Quantity, p.UserID, p.DateCreated, p.DateUpdated); err != nil {
		return nil, errors.Wrap(err, "inserting product")
	}
	return &p, nil
}

func Update(ctx context.Context, db *sqlx.DB, user auth.Claims, id string, update UpdateProduct, now time.Time) error {

	p, err := Retrieve(ctx, db, id)

	if err != nil {
		log.Fatal("error updated")
		return err
	}

	if p.UserID != user.Subject {
		return ErrForbidden
	}

	if update.Name != nil {
		p.Name = *update.Name
	}

	if update.Cost != nil {
		p.Cost = *update.Cost
	}

	if update.Quantity != nil {
		p.Quantity = *update.Quantity
	}

	p.DateUpdated = now

	const q = `UPDATE products SET 
                    "name" = $2, "cost" = $3,
					"quantity" = $4,
                    "date_updated" = $5
					WHERE product_id = $1`

	_, err = db.ExecContext(ctx, q, id,
		p.Name, p.Cost,
		p.Quantity, p.DateUpdated,
	)
	if err != nil {
		return errors.Wrap(err, "updating product")
	}

	return nil
}

func Delete(ctx context.Context, db *sqlx.DB, id string) error {

	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidId
	}
	const q = `DELETE FROM products where product_id = $1`

	_, err := db.ExecContext(ctx, q, id)

	if err != nil {
		return errors.Wrap(err, "deleting product")
	}
	return nil
}
