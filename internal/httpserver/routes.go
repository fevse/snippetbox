package httpserver

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/fevse/snippetbox/pkg/models/mysql"
)

type Application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func NewApp(errorLog *log.Logger, infoLog *log.Logger, snippets *mysql.SnippetModel,
	templateCache map[string]*template.Template) *Application {
	return &Application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      snippets,
		templateCache: templateCache,
	}
}

func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(neuteredFileSistem{http.Dir("./ui/static")})

	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	return mux
}

type neuteredFileSistem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSistem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}
