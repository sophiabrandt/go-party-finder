package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sophiabrandt/go-party-finder/internal/config"
	"github.com/sophiabrandt/go-party-finder/internal/models"
)

var (
	// functions holds all functions that are available in the templates.
	functions = template.FuncMap{}
	conf      *config.Conf
)

// NewTemplates sets the config for the template package.
func NewTemplates(c *config.Conf) {
	conf = c
}

// Respond renders templates using html/template.
func Respond(ctx context.Context, w http.ResponseWriter, tmpl string, data interface{}, statusCode int) error {
	// Set the status code for the request logger middleware.
	// If the context is missing this value, request the service
	// to be shutdown gracefully.
	v, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		return NewShutdownError("web value missing from context")
	}
	v.StatusCode = statusCode

	// write to buffer with HTML tempplate
	if td, ok := data.(*models.TemplateData); ok {

		// setup template Cache
		var tc map[string]*template.Template

		if conf.App.UseCache {
			// get the template cache from the app config
			tc = conf.App.TemplateCache
		} else {
			// this is just used for testing, so that we rebuild
			// the cache on every request
			tc, _ = CreateTemplateCache()
		}

		t, ok := tc[tmpl]
		if !ok {
			return errors.New("can't get template from cache")
		}

		// add secure headers
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Write the status code to the response.
		w.WriteHeader(statusCode)

		buf := new(bytes.Buffer)

		err := t.Execute(buf, td)
		if err != nil {
			log.Fatal(err)
		}
		_, err = buf.WriteTo(w)
		if err != nil {
			fmt.Println("Error writing template to browser", err)
			return err
		}

	} else {
		// If input data is not template data, try marshalling to json
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}

		// Set the content type and headers once we know marshaling has succeeded.
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")

		// Write the status code to the response.
		w.WriteHeader(statusCode)

		// Send the result back to the client.
		if _, err := w.Write(jsonData); err != nil {
			return err
		}
	}

	if err, ok := data.(ErrorResponse); ok {
		http.Error(w, err.Error, http.StatusInternalServerError)
		return nil
	}

	return nil
}

// CreateTemplateCache creates a template cache as a map.
func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", conf.App.TemplateLocation))
	if err != nil {
		return cache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", conf.App.TemplateLocation))
		if err != nil {
			return cache, err
		}

		ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.partial.tmpl", conf.App.TemplateLocation))
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
