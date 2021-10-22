package main

import (
	"path/filepath"
	"text/template"

	"github.com/BunnyTheLifeguard/snipsnip/pkg/models"
)

type templateData struct {
	Snip  *models.Snip
	Snips []*models.Snip
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
		ts, err := template.ParseFiles(page)
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
