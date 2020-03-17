package handlers

import (
	"github.com/harrisonbrock/gargesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func API(logger *log.Logger, db *sqlx.DB) http.Handler {

	app := web.NewApp(logger)
	p := Products{
		DB:  db,
		Log: logger,
	}

	app.Handle(http.MethodGet, "/v1/products", p.List)
	app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve)
	app.Handle(http.MethodPost, "/v1/products", p.Create)

	return app
}
