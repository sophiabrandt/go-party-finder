package mid

import (
	"context"
	"net/http"

	"github.com/golangcollege/sessions"
	"github.com/sophiabrandt/go-party-finder/internal/web"
)

// Session adds middleware to load and save sessions.
func Session(ses *sessions.Session) web.Middleware {
	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {
		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			// Construct http.Handler with ServeHTTP method for the session middleware.
			s := func(w http.ResponseWriter, r *http.Request) {
				handler.ServeHTTP(w, r)
			}
			ses.Enable(http.HandlerFunc(s))

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
