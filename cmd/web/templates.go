package main

import (
	"path/filepath"
	"text/template"
	"time"

	"github.com/BunnyTheLifeguard/snipsnip/pkg/forms"
	"github.com/BunnyTheLifeguard/snipsnip/pkg/models"
)

type templateData struct {
	AuthenticatedUser string
	CSRFToken         string
	CurrentYear       int
	Flash             string
	Form              *forms.Form
	Snip              *models.Snip
	Snips             []*models.Snip
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 06 15:04 MST")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize new map to act as cache
	cache := map[string]*template.Template{}

	// Get a slice of all 'page' templates
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract filenames from full path
		name := filepath.Base(page)

		// Parse page template file into template set
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add any 'layout' templates to template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add any 'partial' templates to template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add template set to cache using name of page as key
		cache[name] = ts
	}

	return cache, nil
}
