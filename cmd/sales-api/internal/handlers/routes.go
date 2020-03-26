package handlers

import (
	"github.com/harrisonbrock/gargesale/internal/mid"
	"github.com/harrisonbrock/gargesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func API(logger *log.Logger, db *sqlx.DB) http.Handler {

	app := web.NewApp(logger, mid.Logger(logger), mid.Errors(logger), mid.Metrics())

	{
		c := Check{DB: db}
		app.Handle(http.MethodGet, "/v1/health", c.Health)
	}

	{

		p := Products{DB: db, Log: logger}
		// Products.
		app.Handle(http.MethodGet, "/v1/products", p.List)
		app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve)
		app.Handle(http.MethodPost, "/v1/products", p.Create)
		app.Handle(http.MethodPut, "/v1/products/{id}", p.Update)
		app.Handle(http.MethodDelete, "/v1/products/{id}", p.Delete)

		// Sales.
		app.Handle(http.MethodPost, "/v1/products/{id}/sales", p.AddSale)
		app.Handle(http.MethodGet, "/v1/products/{id}/sales", p.ListSales)
	}

	return app
}
