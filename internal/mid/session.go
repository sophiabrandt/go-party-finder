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
			// How to enable this?
			// cannot use handler (type web.Handler) as type http.Handler in argument to ses.Enable:
			// web.Handler does not implement http.Handler (missing ServeHTTP method)
			ses.Enable(handler)

			// Call the next handler and set its return value in the err variable.
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
