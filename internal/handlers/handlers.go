package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/sophiabrandt/go-party-finder/internal/mid"
	"github.com/sophiabrandt/go-party-finder/internal/web"
)

// Router  creates a new http.Handler with all routes.
func Router(build string, shutdown chan os.Signal, staticFilesDir string, log *log.Logger) http.Handler {
	// Creates a new web application with all routes and middleware.
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	// Register debug check endpoints.
	cg := checkGroup{
		build: build,
	}
	app.HandleDebug(http.MethodGet, "/readiness", cg.readiness)

	// index route
	app.Handle(http.MethodGet, "/", Home)

	// static file server
	filesDir := http.Dir(staticFilesDir)
	app.FileServer("/static", web.NeuteredFileSystem{filesDir})

	return app
}
