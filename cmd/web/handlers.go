package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	html, err := template.ParseFiles(files...)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = html.ExecuteTemplate(w, "base", nil)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Go can't distinguish JSON from plain text, so it'll always be detected text/plain
func getSnippet(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	snippetId, err := strconv.Atoi(queryParams.Get("id"))

	if err != nil || snippetId < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Snippet with ID %d", snippetId)
}

/**
* Go set 3 system generated headers (Date, Content-Length and Content-Type)
* It attempts to set the correct content type by sniffing the response body with http.DetectContentType function
* Fallback if can't be guessed is application/octet-stream
 */
func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// call writeHeader otherwise write will send a 200 OK
		w.Header().Set("Allow", "POST")

		// helper that calls writeHeader & write in the background
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a snippet"))
}
