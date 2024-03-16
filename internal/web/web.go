package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/notes"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/templ/components"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/templ/page"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/templ/view"
)

func getNotesFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("GetNotes")
	page.Index(view.NotesView(notes.GetAll())).Render(r.Context(), w)
}

func addNoteModalFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("AddNoteModal")
	components.ModalAddNote().Render(r.Context(), w)
}

func addNoteFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("AddNote")

	if r.FormValue("title") == "" {
		w.WriteHeader(http.StatusBadRequest)
		components.ErrorMsg("Title is empty").Render(r.Context(), w)
		return
	}

	if r.FormValue("body") == "" {
		w.WriteHeader(http.StatusBadRequest)
		components.ErrorMsg("Body is empty").Render(r.Context(), w)
		return
	}

	notes.Add(notes.CreateNote{
		Title: r.FormValue("title"),
		Body:  r.FormValue("body"),
	})

	time.Sleep(1 * time.Second)
	view.NotesView(notes.GetAll()).Render(r.Context(), w)
}

func editNoteModalFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("GetEditModal " + r.PathValue("id"))

	note, err := notes.GetNoteByID(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		components.ErrorMsg(err.Error()).Render(r.Context(), w)
		return
	}

	components.ModalEditNote(note).Render(r.Context(), w)
}

func editNoteFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("EditNote")

	if r.FormValue("title") == "" {
		w.WriteHeader(http.StatusBadRequest)
		components.ErrorMsg("Title is empty").Render(r.Context(), w)
		return
	}

	if r.FormValue("body") == "" {
		w.WriteHeader(http.StatusBadRequest)
		components.ErrorMsg("Body is empty").Render(r.Context(), w)
		return
	}

	note, err := notes.GetNoteByID(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		components.ErrorMsg(err.Error()).Render(r.Context(), w)
		return
	}

	note.Title = r.FormValue("title")
	note.Body = r.FormValue("body")
	notes.Update(note)

	time.Sleep(1 * time.Second)
	view.NotesView(notes.GetAll()).Render(r.Context(), w)
}

func deleteNoteFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteNote")
	if err := notes.Delete(r.PathValue("id")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		components.ErrorMsg(err.Error()).Render(r.Context(), w)
		return
	}
	view.NotesView(notes.GetAll()).Render(r.Context(), w)
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

	mux.HandleFunc("/add", addNoteModalFunc)
	mux.HandleFunc("POST /notes", addNoteFunc)

	mux.HandleFunc("/edit/{id}", editNoteModalFunc)
	mux.HandleFunc("PUT /note/{id}", editNoteFunc)
	mux.HandleFunc("DELETE /note/{id}", deleteNoteFunc)

	// Use http.FS to create a file system from the embedded files
	fileSystem, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(fileSystem))))

	log.Fatal(http.ListenAndServe(":8080", mux))
}