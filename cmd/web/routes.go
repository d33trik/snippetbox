package main

import (
	"net/http"

	"codeberg.org/d33trik/snippetbox/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	dynamic := alice.New(app.sessionManager.LoadAndSave, preventCSRF, app.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))

	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))

	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/save", dynamic.ThenFunc(app.userSave))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/authenticate", dynamic.ThenFunc(app.userAuthenticate))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/save", protected.ThenFunc(app.snippetSave))

	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogout))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
