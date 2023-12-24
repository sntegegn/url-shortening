package main

import (
	"bytes"
	"html/template"
	"math/rand"
	"net/http"
	"time"

	"github.com/sntegegn/url-shortening/ui"
)

func (app *application) sendEmail() {
	err := app.mailer.Send("john@example.com", "email.tmpl.html", nil)
	if err != nil {
		app.logger.Error(err.Error())
	}
}

func (app *application) generateShortKey(url string) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, msg string) {
	app.logger.Error(msg)
	http.Error(w, msg, http.StatusInternalServerError)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, statusCode int, page, templateName string, data any) {
	ts, err := template.New("form").ParseFS(ui.Files, "html/"+page)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	buf := new(bytes.Buffer)

	err = ts.ExecuteTemplate(buf, templateName, data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(statusCode)

	buf.WriteTo(w)

}
