package handlers

import (
	"context"
	"net/http"
	"os"

	"github.com/sophiabrandt/go-party-finder/foundation/database"
	"github.com/sophiabrandt/go-party-finder/foundation/web"

	"github.com/jmoiron/sqlx"
)

type checkGroup struct {
	build string
	db    *sqlx.DB
}

// readiness checks if the database is ready and if not will return a 500 status.
// Do not respond by just returning an error because further up in the call
// stack it will interpret that as a non-trusted error.
func (cg checkGroup) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := "ok"
	statusCode := http.StatusOK

	if err := database.StatusCheck(ctx, cg.db); err != nil {
		status = "db not ready"
		statusCode = http.StatusInternalServerError
	}

	health := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	return web.Respond(ctx, w, r, "", health, statusCode)
}

// liveness returns simple status info if the service is alive.
func (cg checkGroup) liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	info := struct {
		Status string `json:"status,omitempty"`
		Build  string `json:"build,omitempty"`
		Host   string `json:"host,omitempty"`
	}{
		Status: "up",
		Build:  cg.build,
		Host:   host,
	}

	return web.Respond(ctx, w, r, "", info, http.StatusOK)
}
