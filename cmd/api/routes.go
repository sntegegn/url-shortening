package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) route() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/shorten", app.shortenURL)
	router.HandlerFunc(http.MethodGet, "/shorten", app.shortenURLForm)
	router.HandlerFunc(http.MethodGet, "/s/:key", app.expandURL)

	return router
}
