package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *App) routes() http.Handler {
	mux := httprouter.New()

	mux.HandlerFunc("POST", "/api/create", app.handleShortUrlCreation)
	mux.HandlerFunc("GET", "/*path", app.handleShortUrlRedirect)

	return mux
}
