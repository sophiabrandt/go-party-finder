package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/sophiabrandt/go-party-finder/internal/mid"
	"github.com/sophiabrandt/go-party-finder/internal/web"
)

// Router  creates a new http.Handler with all routes.
func Router(build string, shutdown chan os.Signal, log *log.Logger) http.Handler {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	app.Handle(http.MethodGet, "/", Home)

	// static file server
	filesDir := http.Dir("./ui/static")
	app.FileServer("/static", filesDir)

	return app
}
