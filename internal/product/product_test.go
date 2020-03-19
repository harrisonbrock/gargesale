package product_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/harrisonbrock/gargesale/internal/product"
	"github.com/harrisonbrock/gargesale/internal/schema"
	"github.com/harrisonbrock/gargesale/internal/tests"
	"testing"
	"time"
)

func TestProducts(t *testing.T) {

	db, cleanup := tests.NewUnit(t)
	defer cleanup()

	np := product.NewProduct{
		Name:     "Toy Gun",
		Cost:     25,
		Quantity: 10,
	}

	now := time.Date(2019, time.July, 1, 0, 0, 0, 0, time.UTC)
	p, err := product.Create(db, np, now)

	if err != nil {
		t.Fatalf("could not create product: %v", err)
	}

	saved, err := product.Retrieve(db, p.ID)

	if err != nil {
		t.Fatalf("could not retrieve product: %v", err)
	}

	if diff := cmp.Diff(p, saved); diff != "" {
		t.Fatalf("saved product did not match created: see diff:\n%s", diff)
	}

}

func TestList(t *testing.T) {
	db, cleanup := tests.NewUnit(t)
	defer cleanup()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	ps, err := product.List(db)

	if err != nil {
		t.Fatalf("listing products: %s", err)
	}

	if exp, got := 2, len(ps); exp != got {
		t.Fatalf("expected product list size %v, got %v", exp, got)
	}
}
