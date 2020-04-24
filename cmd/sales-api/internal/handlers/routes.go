package handlers

import (
	"github.com/harrisonbrock/gargesale/internal/mid"
	"github.com/harrisonbrock/gargesale/internal/platform/auth"
	"github.com/harrisonbrock/gargesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func API(logger *log.Logger, db *sqlx.DB, authenticator *auth.Authenticator) http.Handler {

	app := web.NewApp(logger, mid.Logger(logger), mid.Errors(logger), mid.Metrics())

	{
		c := Check{DB: db}
		app.Handle(http.MethodGet, "/v1/health", c.Health)
	}

	{

		p := Products{DB: db, Log: logger}
		// Products.
		app.Handle(http.MethodGet, "/v1/products", p.List, mid.Authenticate(authenticator))
		app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve, mid.Authenticate(authenticator))
		app.Handle(http.MethodPost, "/v1/products", p.Create, mid.Authenticate(authenticator))
		app.Handle(http.MethodPut, "/v1/products/{id}", p.Update, mid.Authenticate(authenticator))
		app.Handle(http.MethodDelete, "/v1/products/{id}", p.Delete, mid.Authenticate(authenticator))

		// Sales.
		app.Handle(http.MethodPost, "/v1/products/{id}/sales", p.AddSale, mid.Authenticate(authenticator))
		app.Handle(http.MethodGet, "/v1/products/{id}/sales", p.ListSales, mid.Authenticate(authenticator))
	}

	{
		u := Users{DB: db, authenticator: authenticator}
		app.Handle(http.MethodGet, "/v1/users/token", u.Token)
	}

	return app
}
