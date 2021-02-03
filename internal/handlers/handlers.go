package handlers

import (
	"net/http"

	"github.com/sophiabrandt/go-party-finder/internal/models"
	"github.com/sophiabrandt/go-party-finder/internal/render"
)

func Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}
