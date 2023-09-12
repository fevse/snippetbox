package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Showing a note"))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Form for creating a new note"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Println("Running a web-server on http://127.0.0.1.4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
