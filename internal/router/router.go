package router

import (
	"log"
	"net/http"

	"github.com/sophiabrandt/go-party-finder/internal/handlers"
	"github.com/sophiabrandt/go-party-finder/internal/web"
)

// New creates a new http.Handler with all routes.
func New(log *log.Logger) http.Handler {
	app := web.NewApp(log)

	app.Handle(http.MethodGet, "/", handlers.Home)

	// static file server
	filesDir := http.Dir("./ui/static")
	app.FileServer("/static", filesDir)

	return app
}
