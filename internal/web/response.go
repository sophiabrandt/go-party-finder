package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/sophiabrandt/go-party-finder/internal/config"
	td "github.com/sophiabrandt/go-party-finder/internal/data"
)

var (
	conf *config.Conf

	// functions holds all functions that are available in the templates.
	functions = template.FuncMap{
		"humanDate": HumanDate,
	}
)

// NewTemplates sets the config for the template package.
func NewTemplates(c *config.Conf) {
	conf = c
}

// HumanDate returns time in YYYY-MM-DD format
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// addDefaultData adds data for all templates
func addDefaultData(dt *td.TemplateData) *td.TemplateData {
	if dt == nil {
		dt = &td.TemplateData{}
	}
	dt.CurrentYear = time.Now().Year()

	return dt
}

// Respond renders templates using html/template.
func Respond(ctx context.Context, w http.ResponseWriter, tmpl string, data interface{}, statusCode int) error {
	// add secure headers
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Frame-Options", "deny")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Set the status code for the request logger middleware.
	// If the context is missing this value, request the service
	// to be shutdown gracefully.
	v, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		return NewShutdownError("web value missing from context")
	}
	v.StatusCode = statusCode

	switch d := data.(type) {
	case *td.TemplateData:
		return respondWithTemplate(ctx, w, tmpl, d, statusCode)
	case ErrorResponse:
		return respondWithError(w, d, statusCode)
	default:
		return respondWithJson(w, d, statusCode)
	}

	return nil
}

// respondWithTemplate assembles the HTML template and renders a response to the client.
func respondWithTemplate(ctx context.Context, w http.ResponseWriter, tmpl string, data *td.TemplateData, statusCode int) error {
	// setup template Cache
	var tc map[string]*template.Template

	if conf.Web.UseCache {
		// get the template cache from the app config
		tc = conf.Web.TemplateCache
	} else {
		// this is just used for testing, so that we rebuild
		// the cache on every request
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		return NewShutdownError("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	err := t.Execute(buf, addDefaultData(data))
	if err != nil {
		return errors.Wrap(err, "cannot parse template")
	}

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	_, err = buf.WriteTo(w)
	if err != nil {
		return errors.Wrap(err, "Error writing template to browser")
	}
	return nil
}

// respondWithJson marshalls data into json and returns it to the client.
func respondWithJson(w http.ResponseWriter, data interface{}, statusCode int) error {
	// If input data is not template data, try marshalling to json
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}
	return nil
}

// respondWithError renders an HTTP error to the client.
func respondWithError(w http.ResponseWriter, data ErrorResponse, statusCode int) error {
	http.Error(w, data.Error, http.StatusInternalServerError)
	return nil
}

// CreateTemplateCache creates a template cache as a map.
func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", conf.Web.TemplateLocation))
	if err != nil {
		return cache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", conf.Web.TemplateLocation))
		if err != nil {
			return cache, err
		}

		ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.partial.tmpl", conf.Web.TemplateLocation))
		if err != nil {
			return cache, err
		}

		cache[name] = ts
	}

	return cache, nil
}

// RespondError sends an error response back to the client.
func RespondError(ctx context.Context, w http.ResponseWriter, err error) error {
	// If the error was of the type *Error, the handler has
	// a specific status code and error to return.
	if webErr, ok := errors.Cause(err).(*Error); ok {
		er := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		if err := Respond(ctx, w, "", er, webErr.Status); err != nil {
			return err
		}
		return nil
	}

	// If not, the handler sent any arbitrary error value so use 500.
	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	if err := Respond(ctx, w, "", er, http.StatusInternalServerError); err != nil {
		return err
	}

	return nil
}
