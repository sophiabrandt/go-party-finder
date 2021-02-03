package handlers

import (
	"context"
	"net/http"

	"github.com/sophiabrandt/go-party-finder/internal/models"
	"github.com/sophiabrandt/go-party-finder/internal/web"
)

func Home(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	web.Respond(ctx, w, "home.page.tmpl", &models.TemplateData{}, http.StatusOK)
	return nil
}
