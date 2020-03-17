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
func (p *Products) List(w http.ResponseWriter, r *http.Request) {

	list, err := product.List(p.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error querying data source", err)
		return
	}

	if err := web.Response(w, list, http.StatusOK); err != nil {
		p.Log.Println("error responding", err)
		return
	}
}

func (p *Products) Retrieve(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	prod, err := product.Retrieve(p.DB, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error querying data source", err)
		return
	}

	if err := web.Response(w, prod, http.StatusOK); err != nil {
		p.Log.Println("error responding", err)
		return
	}
}

// Create decode json document from a POST Request
func (p *Products) Create(w http.ResponseWriter, r *http.Request) {

	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		p.Log.Println(err)
		return
	}

	prod, err := product.Create(p.DB, np, time.Now())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error querying data source", err)
		return
	}

	if err := web.Response(w, prod, http.StatusCreated); err != nil {
		p.Log.Println("error responding", err)
		return
	}
}
