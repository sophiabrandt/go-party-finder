package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/sophiabrandt/go-party-finder/internal/config"
	"github.com/sophiabrandt/go-party-finder/internal/data/party"
	"github.com/sophiabrandt/go-party-finder/internal/mid"
	"github.com/sophiabrandt/go-party-finder/internal/web"
)

// Router  creates a new http.Handler with all routes.
func Router(build string, shutdown chan os.Signal, lctx *config.LocalContext, staticFilesDir string, log *log.Logger, db *sqlx.DB) http.Handler {
	// Creates a new web application with all routes and middleware.
	app := web.NewApp(shutdown, lctx, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log), mid.Session(lctx.Session))

	// Register debug check endpoints.
	cg := checkGroup{
		build: build,
		db:    db,
	}
	app.HandleDebug(http.MethodGet, "/readiness", cg.readiness)
	app.HandleDebug(http.MethodGet, "/liveness", cg.liveness)

	// index route and parties routes
	pg := partyGroup{
		party: party.New(log, db),
	}
	app.Handle(http.MethodGet, "/", pg.query, mid.Session(lctx.Session))
	app.Handle(http.MethodGet, "/parties/{page}/{rows}", pg.query, mid.Session(lctx.Session))
	app.Handle(http.MethodGet, "/parties/{id}", pg.queryByID, mid.Session(lctx.Session))
	app.Handle(http.MethodGet, "/parties/create", pg.createForm, mid.Session(lctx.Session))
	app.Handle(http.MethodPost, "/parties/create", pg.create, mid.Session(lctx.Session))

	// static file server
	filesDir := http.Dir(staticFilesDir)
	app.FileServer("/static", web.NeuteredFileSystem{filesDir})

	return app
}
