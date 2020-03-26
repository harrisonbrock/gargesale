package handlers

import (
	"context"
	"github.com/harrisonbrock/gargesale/internal/platform/database"
	"github.com/harrisonbrock/gargesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Check struct {
	DB *sqlx.DB
}

func (c *Check) Health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	var health struct {
		Status string `json:"status"`
	}

	if err := database.StatusCheck(ctx, c.DB); err != nil {
		health.Status = "database not ready"
		return web.Respond(ctx, w, health, http.StatusInternalServerError)
	}
	health.Status = "database ready"
	return web.Respond(ctx, w, health, http.StatusOK)
}
