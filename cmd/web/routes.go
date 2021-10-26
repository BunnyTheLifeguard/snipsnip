package main

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := chi.NewRouter()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home).ServeHTTP)
	mux.Get("/snip/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createSnipForm).ServeHTTP)
	mux.Post("/snip/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createSnip).ServeHTTP)
	mux.Get("/snip/{id}", dynamicMiddleware.ThenFunc(app.showSnip).ServeHTTP)

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm).ServeHTTP)
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser).ServeHTTP)
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm).ServeHTTP)
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser).ServeHTTP)
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.logoutUser).ServeHTTP)

	filesDir := http.Dir("../../ui/static")
	fileServer(mux, "/static", filesDir)

	return standardMiddleware.Then(mux)
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
