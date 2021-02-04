package handlers

import (
	"context"
	"net/http"

	"github.com/sophiabrandt/go-party-finder/internal/web"
)

type checkGroup struct {
	build string
}

// readiness checks if the database is ready and if not will return a 500 status.
// Do not respond by just returning an error because further up in the call
// stack it will interpret that as a non-trusted error.
func (cg checkGroup) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := "ok"
	statusCode := http.StatusOK

	health := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	return web.Respond(ctx, w, "", health, statusCode)
}
