package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sntegegn/url-shortening/internal/models"
	"github.com/sntegegn/url-shortening/internal/validator"
)

type urlForm struct {
	LongURL string `form:"longURL"`
	validator.Validator
}

type templateData struct {
	LongURL  string
	ShortURL string
	Errors   map[string]string
}

func (app *application) shortenURL(w http.ResponseWriter, r *http.Request) {
	ctx, span := app.tracer.Start(r.Context(), "shortenURL_post")
	defer span.End()

	form := new(urlForm)

	err := app.decodePostForm(r, form)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	form.CheckField(validator.NotBlank(form.LongURL), "longURL", "LongURL cannot be blank")
	form.CheckField(validator.Matches(form.LongURL, validator.URLRX), "longURL", "LongURL must be a valid URL")

	if !form.Valid() {
		data := templateData{
			Errors: form.FieldError,
		}
		app.render(w, r, http.StatusBadRequest, "form.tmpl.html", "urlForm", data)
		return
	}

	shortKey := app.generateShortKey(ctx)

	err = app.URLModel.Insert(shortKey, form.LongURL)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	shortURL := fmt.Sprintf("http://localhost:4000/s/%s", shortKey)
	data := templateData{
		LongURL:  form.LongURL,
		ShortURL: shortURL,
	}

	app.render(w, r, http.StatusOK, "form.tmpl.html", "urlForm", data)

}

func (app *application) shortenURLForm(w http.ResponseWriter, r *http.Request) {
	_, span := app.tracer.Start(r.Context(), "shortenURL_get")
	defer span.End()

	app.render(w, r, http.StatusOK, "form.tmpl.html", "urlForm", nil)
}

func (app *application) expandURL(w http.ResponseWriter, r *http.Request) {
	_, span := app.tracer.Start(r.Context(), "expandURL")
	defer span.End()

	params := httprouter.ParamsFromContext(r.Context())
	shortKey := params.ByName("key")

	longURL, err := app.URLModel.Get(shortKey)
	if err != nil {
		expandCount.WithLabelValues("failed").Inc()
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.badRequestError(w, r, err)
		default:
			app.serverError(w, r, err)
		}
	}
	expandCount.WithLabelValues("success").Inc()

	http.Redirect(w, r, longURL, http.StatusSeeOther)
}
