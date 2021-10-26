package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CSRFToken = nosurf.Token(r)
	td.AuthenticatedUser = app.authenticatedUser(r)
	td.CurrentYear = time.Now().Year()
	td.Flash = app.session.PopString(r, "flash")

	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	// Initialize buffer for catching runtime errors
	buf := new(bytes.Buffer)

	// Write template to buffer instead of ResponseWriter
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
	}

	// Write contents of buffer to ResponseWriter if no error present
	buf.WriteTo(w)
}

func (app application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) authenticatedUser(r *http.Request) string {
	return app.session.GetString(r, "userID")
}
