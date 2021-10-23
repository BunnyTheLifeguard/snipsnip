package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/BunnyTheLifeguard/snipsnip/pkg/models"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snips.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snips: s,
	})
}

func (app *application) showSnip(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
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

	app.render(w, r, "show.page.tmpl", &templateData{
		Snip: s,
	})
}

func (app *application) createSnipForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}

func (app *application) createSnip(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	days := r.PostForm.Get("expires")
	daysInt, _ := strconv.Atoi(days)
	expires := time.Now().Add(time.Hour * 24 * time.Duration(daysInt))
	created := time.Now()

	id, err := app.snips.Insert(title, content, created, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	oid := id.(primitive.ObjectID).Hex()

	http.Redirect(w, r, fmt.Sprintf("/snip/%s", oid), http.StatusSeeOther)
}
