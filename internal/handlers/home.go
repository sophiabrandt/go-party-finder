package handlers

import (
	"context"
	"net/http"

	"github.com/sophiabrandt/go-party-finder/internal/data"
	"github.com/sophiabrandt/go-party-finder/internal/web"
)

func Home(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	web.Respond(ctx, w, "home.page.tmpl", &data.TemplateData{}, http.StatusOK)
	return nil
}
