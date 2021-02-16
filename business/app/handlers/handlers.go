package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/sophiabrandt/go-party-finder/business/app/apc"
	"github.com/sophiabrandt/go-party-finder/business/data/party"
	"github.com/sophiabrandt/go-party-finder/business/mid"
	"github.com/sophiabrandt/go-party-finder/foundation/web"
)

// Router creates a new http.Handler with all routes.
func Router(build string, shutdown chan os.Signal, apc *apc.AppContext, staticFilesDir string, log *log.Logger, db *sqlx.DB) http.Handler {
	// Creates a new web application with all routes and middleware.
	app := web.NewApp(shutdown, apc, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	// Register debug check endpoints.
	cg := checkGroup{
		build: build,
		db:    db,
	}
	app.HandleDebug(http.MethodGet, "/readiness", web.Handler{cg.readiness})
	app.HandleDebug(http.MethodGet, "/liveness", web.Handler{cg.liveness})

	// index route and parties routes
	pg := partyGroup{
		party: party.New(log, db),
	}
	app.Handle(http.MethodGet, "/", web.Handler{pg.query})
	app.Handle(http.MethodGet, "/parties/{page}/{rows}", web.Handler{pg.query})
	app.Handle(http.MethodGet, "/parties/{id}", web.Handler{pg.queryByID})
	app.Handle(http.MethodGet, "/parties/create", web.Handler{pg.createForm})
	app.Handle(http.MethodPost, "/parties/create", web.Handler{pg.create})

	// static file server
	filesDir := http.Dir(staticFilesDir)
	app.FileServer("/static", web.NeuteredFileSystem{filesDir})

	// TODO: find a better way to add middleware with http.Handler signature to
	// individual routes
	return apc.Session.Enable(app)
}
