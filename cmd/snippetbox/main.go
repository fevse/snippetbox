package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	h "github.com/fevse/snippetbox/internal/httpserver"
	"github.com/fevse/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	addr := flag.String("addr", ":4000", "Net address HTTP")
	dsn := flag.String("dsn", "web:qseft135@/snippetbox?parseTime=true", "MySQL Data Source Name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := h.NewTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := h.NewApp(errorLog, infoLog, &mysql.SnippetModel{DB: db}, templateCache)

	srv := &http.Server{
		Addr:              *addr,
		ErrorLog:          errorLog,
		Handler:           app.Routes(),
		ReadHeaderTimeout: 2 * time.Second,
	}

	infoLog.Printf("Running a web-server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
