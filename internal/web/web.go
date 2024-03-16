package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/notes"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/templ/page"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/templ/view"
)

func getNotesFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("GetNotes")
	page.Index(view.NotesView(notes.GetAll())).Render(r.Context(), w)
}

//go:embed static/*
var staticFiles embed.FS

func Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the requested URL is one of the defined handlers
		// If not, redirect to the custom 404 page
		_, pattern := mux.Handler(r)
		if r.URL.Path != pattern {
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		// Serve the index handler
		getNotesFunc(w, r)
	})

	mux.Handle("/404", templ.Handler(page.Index(view.NotFoundComponent()),
		templ.WithStatus(http.StatusNotFound)))

	mux.HandleFunc("/notes", getNotesFunc)

	// Use http.FS to create a file system from the embedded files
	fileSystem, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(fileSystem))))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
