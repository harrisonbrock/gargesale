package handlers

import (
	"encoding/json"
	"github.com/harrisonbrock/gargesale/internal/product"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

// Products has handler methods for dealing with Products.
type Products struct {
	DB *sqlx.DB
}

// ListProduct is a basic HTTP Handler.
func (p *Products) List(w http.ResponseWriter, r *http.Request) {

	list, err := product.List(p.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error querying data source", err)
		return
	}

	data, err := json.Marshal(list)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error marshalling", err)
		return
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Println("error writing", err)
	}

}
