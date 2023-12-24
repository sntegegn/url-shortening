package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type urlForm struct {
	LongURL string `form:"longURL"`
}

type templateData struct {
	LongURL  string
	ShortURL string
}

func (app *application) shortenURL(w http.ResponseWriter, r *http.Request) {
	form := new(urlForm)

	err := app.decodePostForm(r, form)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	shortKey := app.generateShortKey(form.LongURL)
	shortURL := fmt.Sprintf("http://localhost:4000/s/%s", shortKey)
	app.urls[shortKey] = form.LongURL
	data := templateData{
		LongURL:  form.LongURL,
		ShortURL: shortURL,
	}
	app.render(w, r, http.StatusOK, "result.tmpl.html", "result", data)

}

func (app *application) shortenURLForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "form.tmpl.html", "urlForm", nil)
}

func (app *application) expandURL(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	shortKey := params.ByName("key")
	longURL, ok := app.urls[shortKey]
	if !ok {
		msg := "No longURL corresponding the shortURL found"
		app.badRequestError(w, r, msg)
	}
	http.Redirect(w, r, longURL, http.StatusSeeOther)
}
