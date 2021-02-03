package web

import (
	"net/http"

	"github.com/go-chi/chi"
)

// Param returns the web call parameter from the request.
func Param(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}
