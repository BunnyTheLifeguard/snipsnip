package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/BunnyTheLifeguard/snipsnip/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	s, err := app.snips.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Snips: s}

	files := []string{
		"../../ui/html/home.page.tmpl",
		"../../ui/html/base.layout.tmpl",
		"../../ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) showSnip(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		app.notFound(w)
		return
	}

	s, err := app.snips.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Create instance of a templateData struct holding the snip data
	data := &templateData{Snip: s}

	// Initialize template files
	files := []string{
		"../../ui/html/show.page.tmpl",
		"../../ui/html/base.layout.tmpl",
		"../../ui/html/footer.partial.tmpl",
	}

	// Parse template files
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Execute template files
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) createSnip(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	//TODO: Dummy data - remove later
	title := "3 Snip"
	content := "3 Snip snip!"
	created := time.Now()
	expires := time.Now().Add(time.Hour * 24 * 7)

	id, err := app.snips.Insert(title, content, created, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	oid := id.(primitive.ObjectID).Hex()

	http.Redirect(w, r, fmt.Sprintf("/snip?id=%s", oid), http.StatusSeeOther)
}
