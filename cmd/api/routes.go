package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (app *application) route() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/shorten", app.shortenURL)
	router.HandlerFunc(http.MethodGet, "/shorten", app.shortenURLForm)
	router.HandlerFunc(http.MethodGet, "/s/:key", app.expandURL)
	router.Handler(http.MethodGet, "/metrics", promhttp.Handler())

	return AddMetrics(router)
}
