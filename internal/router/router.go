package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sophiabrandt/go-party-finder/internal/app"
)

// New creates a new http.Handler with all routes.
func New(log *log.Logger) http.Handler {
	app := app.New(log)

	app.Handle(http.MethodGet, "/", Greet)

	return app
}

func Greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
