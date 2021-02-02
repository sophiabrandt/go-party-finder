package app

import (
	"log"
	"net/http"

	"github.com/majidsajadi/sariaf"
)

// App is the entrypoint for the web application.
type App struct {
	mux *sariaf.Router
	log *log.Logger
}

// New creates an App value that handles the routes for the application.
func New(log *log.Logger) *App {
	mux := sariaf.New()

	return &App{
		mux: mux,
		log: log,
	}
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

// Handle is the route handler and wraps the router.
func (a *App) Handle(method string, pattern string, handler http.HandlerFunc) {
	a.mux.Handle(method, pattern, handler)
}
