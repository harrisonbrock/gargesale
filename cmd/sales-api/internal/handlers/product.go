package handlers

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/harrisonbrock/gargesale/internal/platform/auth"
	"github.com/harrisonbrock/gargesale/internal/platform/web"
	"github.com/harrisonbrock/gargesale/internal/product"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

// Products has handler methods for dealing with Products.
type Products struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// ListProduct is a basic HTTP Handler.
func (p *Products) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	list, err := product.List(ctx, p.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error querying data source", err)
		return err
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

func (p *Products) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")

	prod, err := product.Retrieve(ctx, p.DB, id)
	if err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidId:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for product %q", id)
		}
	}

	return web.Respond(ctx, w, prod, http.StatusOK)
}

// Create decode json document from a POST Request
func (p *Products) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)

	if !ok {
		return errors.New("auth claims not in context")
	}

	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		return err
	}

	prod, err := product.Create(ctx, p.DB, claims, np, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, prod, http.StatusCreated)
}

func (p *Products) AddSale(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var ns product.NewSale
	if err := web.Decode(r, &ns); err != nil {
		return errors.Wrap(err, "decoding new sale")
	}

	productID := chi.URLParam(r, "id")

	sale, err := product.AddSale(ctx, p.DB, ns, productID, time.Now())
	if err != nil {
		return errors.Wrap(err, "adding new sale")
	}

	return web.Respond(ctx, w, sale, http.StatusCreated)
}

// ListSales gets all sales for a particular product.
func (p *Products) ListSales(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	list, err := product.ListSales(ctx, p.DB, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

func (p *Products) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	claims, ok := ctx.Value(auth.Key).(auth.Claims)

	if !ok {
		return errors.New("auth claims not in context")
	}

	var update product.UpdateProduct
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding product update")
	}

	if err := product.Update(ctx, p.DB, claims, id, update, time.Now()); err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidId:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating product %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (p *Products) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")

	if err := product.Delete(ctx, p.DB, id); err != nil {
		switch err {
		case product.ErrInvalidId:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting product %q", id)
		}
	}
	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
