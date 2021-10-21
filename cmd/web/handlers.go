package main

import (
	"fmt"
	"net/http"
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

	for _, snip := range s {
		fmt.Fprintf(w, "%v\n", snip)
	}

	// files := []string{
	// 	"../../ui/html/home.page.tmpl",
	// 	"../../ui/html/base.layout.tmpl",
	// 	"../../ui/html/footer.partial.tmpl",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// err = ts.Execute(w, nil)
	// if err != nil {
	// 	app.serverError(w, err)
	// }
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

	fmt.Fprintf(w, "%v", s)
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
