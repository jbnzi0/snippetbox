package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func openDB(datasource string) (*sql.DB, error) {
	db, err := sql.Open("mysql", datasource)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	// this function should be called before using addr
	// go run ./cmd/web -help will list all available cmd line flags available
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile|log.LUTC)

	err := godotenv.Load("./.env")

	if err != nil {
		errorLog.Fatal("Error loading .env file")
	}

	db, err := openDB(os.Getenv("DATABASE_URL"))

	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	infoLog.Printf("Starting server on %s", *addr)

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	// requests are handled concurrently, all incoming requests are served in their own goroutine = higher risk of race conditions when accessing shared resources
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}
