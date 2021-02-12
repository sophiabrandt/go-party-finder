package web

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/sophiabrandt/go-party-finder/internal/apc"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values are stored/retrieved.
const KeyValues ctxKey = 1

// Values represent state for each request.
type Values struct {
	TraceID    string
	Now        time.Time
	AppContext *apc.AppContext
	StatusCode int
}

// registered keeps track of handlers registered to the http default server
// mux. This is a singleton and used by the standard library for metrics
// and profiling. The application may want to add other handlers like
// readiness and liveness to that mux. If this is not tracked, the routes
// could try to be registered more than once, causing a panic.
var registered = make(map[string]bool)

// A Handler is a type that handles an http request within our own little mini
// framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// ServeHTTP is a wrapper to make the Handler compliant with the http.Handler interface.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h(ctx, w, r)
}

// App is the entrypoint for the web application.
type App struct {
	mux      *chi.Mux
	shutdown chan os.Signal
	apc      *apc.AppContext
	mw       []Middleware
}

// NewApp creates an App value that handles the routes for the application.
func NewApp(shutdown chan os.Signal, apc *apc.AppContext, mw ...Middleware) *App {
	mux := chi.NewRouter()

	return &App{
		mux:      mux,
		shutdown: shutdown,
		apc:      apc,
		mw:       mw,
	}
}

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

// HandleDebug sets a handler function for a given HTTP method and path pair
// to the default http package server mux. /debug is added to the path.
func (a *App) HandleDebug(method string, path string, handler Handler, mw ...Middleware) {
	a.handle(true, method, path, handler)
}

// Handle is the route handler and wraps the router.
func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {
	a.handle(false, method, path, handler)
}

// handle provides boilerplate and middleware wrapping.
func (a *App) handle(debug bool, method string, path string, handler Handler, mw ...Middleware) {
	if debug {
		// Track all the handlers that are being registered so we don't have
		// the same handlers registered twice to this singleton.
		if _, exists := registered[method+path]; exists {
			return
		}
		registered[method+path] = true
	}

	// First wrap handler specific middleware around this handler.
	handler = wrapMiddleware(mw, handler)

	// Add the application's general middleware to the handler chain.
	handler = wrapMiddleware(a.mw, handler)

	// The function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request) {
		// Start or expand a distributed trace.
		ctx := r.Context()

		// Set the context with the required values to
		// process the request.
		v := Values{
			TraceID:    uuid.NewString(),
			Now:        time.Now(),
			AppContext: a.apc,
		}
		ctx = context.WithValue(ctx, KeyValues, &v)

		// Call the wrapped handler functions.
		if err := handler(ctx, w, r); err != nil {
			a.SignalShutdown()
			return
		}
	}

	// Add this handler for the specified verb and route.
	if debug {
		f := func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == method:
				h(w, r)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}
		http.DefaultServeMux.HandleFunc("/debug"+path, f)
		return
	}

	a.mux.MethodFunc(method, path, h)
}

// neuteredFileSystem disallows directory listings for a static file server.
type NeuteredFileSystem struct {
	Fs http.FileSystem
}

// Open creates access to the neutered file system.
// https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
func (nfs NeuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.Fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.Fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func (a *App) FileServer(path string, root NeuteredFileSystem) {
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
