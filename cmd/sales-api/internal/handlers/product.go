package handlers

import (
	"github.com/go-chi/chi"
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
func (p *Products) List(w http.ResponseWriter, r *http.Request) error {

	list, err := product.List(r.Context(), p.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error querying data source", err)
		return err
	}

	return web.Response(w, list, http.StatusOK)
}

func (p *Products) Retrieve(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")

	prod, err := product.Retrieve(r.Context(), p.DB, id)
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

	return web.Response(w, prod, http.StatusOK)
}

// Create decode json document from a POST Request
func (p *Products) Create(w http.ResponseWriter, r *http.Request) error {

	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		return err
	}

	prod, err := product.Create(r.Context(), p.DB, np, time.Now())
	if err != nil {
		return err
	}

	return web.Response(w, prod, http.StatusCreated)
}

func (p *Products) AddSale(w http.ResponseWriter, r *http.Request) error {
	var ns product.NewSale
	if err := web.Decode(r, &ns); err != nil {
		return errors.Wrap(err, "decoding new sale")
	}

	productID := chi.URLParam(r, "id")

	sale, err := product.AddSale(r.Context(), p.DB, ns, productID, time.Now())
	if err != nil {
		return errors.Wrap(err, "adding new sale")
	}

	return web.Response(w, sale, http.StatusCreated)
}

// ListSales gets all sales for a particular product.
func (p *Products) ListSales(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	list, err := product.ListSales(r.Context(), p.DB, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Response(w, list, http.StatusOK)
}

func (p *Products) Update(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update product.UpdateProduct
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding product update")
	}

	if err := product.Update(r.Context(), p.db, id, update, time.Now()); err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidId:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating product %q", id)
		}
	}

	return web.Response(w, nil, http.StatusNoContent)
}
