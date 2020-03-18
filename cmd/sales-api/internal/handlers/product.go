package handlers

import (
	"github.com/go-chi/chi"
	"github.com/harrisonbrock/gargesale/internal/platform/web"
	"github.com/harrisonbrock/gargesale/internal/product"
	"github.com/jmoiron/sqlx"
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

	list, err := product.List(p.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error querying data source", err)
		return err
	}

	return web.Response(w, list, http.StatusOK)
}

func (p *Products) Retrieve(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")

	prod, err := product.Retrieve(p.DB, id)
	if err != nil {
		return err
	}

	return web.Response(w, prod, http.StatusOK)
}

// Create decode json document from a POST Request
func (p *Products) Create(w http.ResponseWriter, r *http.Request) error {

	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		return err
	}

	prod, err := product.Create(p.DB, np, time.Now())
	if err != nil {
		return err
	}

	return web.Response(w, prod, http.StatusCreated)
}
