package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/internal/models"
	"strconv"
	"text/template"
)

// we use dependency injection to use custom logs
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.GetSnippets(10)

	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	html, err := template.ParseFiles(files...)

	if err != nil {
		app.serverError(w, err)
		return
	}

	err = html.ExecuteTemplate(w, "base", nil)

	if err != nil {
		app.serverError(w, err)
	}
}

// Go can't distinguish JSON from plain text, so it'll always be detected text/plain
func (app *application) getSnippet(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	id, err := strconv.Atoi(queryParams.Get("id"))

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	fmt.Fprintf(w, "%+v", snippet)
}

/**
* Go set 3 system generated headers (Date, Content-Length and Content-Type)
* It attempts to set the correct content type by sniffing the response body with http.DetectContentType function
* Fallback if can't be guessed is application/octet-stream
 */
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// call writeHeader otherwise write will send a 200 OK
		w.Header().Set("Allow", "POST")

		// helper that calls writeHeader & write in the background
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "0 snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
