package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/BunnyTheLifeguard/snipsnip/pkg/forms"
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
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createSnip(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "1", "3", "7")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
	} else {
		days := form.Get("expires")
		daysInt, _ := strconv.Atoi(days)
		expires := time.Now().Add(time.Hour * 24 * time.Duration(daysInt))
		created := time.Now()

		id, err := app.snips.Insert(form.Get("title"), form.Get("content"), created, expires)
		if err != nil {
			app.serverError(w, err)
			return
		}
		oid := id.(primitive.ObjectID).Hex()

		app.session.Put(r, "flash", "Snip successfully created.")
		http.Redirect(w, r, fmt.Sprintf("/snip/%s", oid), http.StatusSeeOther)
	}
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display user signup form...")
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new user...")
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display user login form...")
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User login...")
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User logout...")
}
