package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	// this function should be called before using addr
	// go run ./cmd/web -help will list all available cmd line flags available
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile|log.LUTC)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.getSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	infoLog.Printf("Starting server on %s", *addr)

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}
	// requests are handled concurrently, all incoming requests are served in their own goroutine = higher risk of race conditions when accessing shared resources
	err := server.ListenAndServe()
	errorLog.Fatal(err)
}
