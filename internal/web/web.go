package web

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

// App is the entrypoint for the web application.
type App struct {
	mux *chi.Mux
	log *log.Logger
}

// NewApp creates an App value that handles the routes for the application.
func NewApp(log *log.Logger) *App {
	mux := chi.NewRouter()

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
	a.mux.Method(method, pattern, handler)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func (a *App) FileServer(path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		a.mux.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	a.mux.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
