package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/sophiabrandt/go-party-finder/internal/config"
	"github.com/sophiabrandt/go-party-finder/internal/models"
)

// functions holds all functions that are available in the templates.
var functions = template.FuncMap{}

var (
	conf            *config.Conf
	pathToTemplates = "./ui/html"
)

// NewTemplates sets the config for the template package.
func NewTemplates(c *config.Conf) {
	conf = c
}

// Template renders templates using html/template.
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
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

	return nil
}

// CreateTemplateCache creates a template cache as a map.
func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return cache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return cache, err
		}

		ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.partial.tmpl", pathToTemplates))
		if err != nil {
			return cache, err
		}

		cache[name] = ts
	}

	return cache, nil
}
