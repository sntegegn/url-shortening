package main

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/sntegegn/url-shortening/ui"
)

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
