package web

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/notes"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/templ/components"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/templ/page"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/templ/view"
)

func handleRenderError(err error) {
	if err != nil {
		fmt.Println("Render error: ", err)
	}
}

// Set header and render error message
func sendErrorMsg(w http.ResponseWriter, r *http.Request, errorMsg string) {
	w.WriteHeader(http.StatusBadRequest)
	handleRenderError(components.ErrorMsg(errorMsg).Render(r.Context(), w))
}

func notesFunc(w http.ResponseWriter, r *http.Request) {
	handleRenderError(page.Index(view.NotesView(notes.GetLatestNotes())).Render(r.Context(), w))
}

func moreNotesFunc(w http.ResponseWriter, r *http.Request) {
	noteID := -1
	if p := r.URL.Query().Get("note"); p != "" {
		if parsedNoteID, err := strconv.Atoi(p); err == nil {
			noteID = parsedNoteID
		}
	}

	if noteID == -1 {
		sendErrorMsg(w, r, "Note is empty")
		return
	}

	notesOnPage := notes.GetNextNotes(noteID)

	time.Sleep(250 * time.Millisecond)
	handleRenderError(components.NotesList(notesOnPage).Render(r.Context(), w))
}

func addNoteModalFunc(w http.ResponseWriter, r *http.Request) {
	handleRenderError(components.ModalAddNote().Render(r.Context(), w))
}

func addNoteFunc(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("title") == "" {
		sendErrorMsg(w, r, "Title is empty")
		return
	}

	if r.FormValue("body") == "" {
		sendErrorMsg(w, r, "Body is empty")
		return
	}

	note := notes.Add(notes.CreateNote{
		Title: r.FormValue("title"),
		Body:  r.FormValue("body"),
	})

	time.Sleep(250 * time.Millisecond)

	handleRenderError(components.NoteItem(note).Render(r.Context(), w))
}

func editNoteModalFunc(w http.ResponseWriter, r *http.Request) {
	note, err := notes.GetNoteByID(r.PathValue("id"))
	if err != nil {
		sendErrorMsg(w, r, err.Error())
		return
	}

	handleRenderError(components.ModalEditNote(note).Render(r.Context(), w))
}

func editNoteFunc(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("title") == "" {
		sendErrorMsg(w, r, "Title is empty")
		return
	}

	if r.FormValue("body") == "" {
		sendErrorMsg(w, r, "Body is empty")
		return
	}

	note, err := notes.GetNoteByID(r.PathValue("id"))
	if err != nil {
		sendErrorMsg(w, r, err.Error())
		return
	}

	note.Title = r.FormValue("title")
	note.Body = r.FormValue("body")
	notes.Update(note)

	time.Sleep(250 * time.Millisecond)
	handleRenderError(components.NoteItem(note).Render(r.Context(), w))
}

func deleteNoteFunc(w http.ResponseWriter, r *http.Request) {
	if err := notes.Delete(r.PathValue("id")); err != nil {
		sendErrorMsg(w, r, err.Error())
	}
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
		notesFunc(w, r)
	})

	mux.Handle("/404", templ.Handler(page.Index(view.NotFoundComponent()),
		templ.WithStatus(http.StatusNotFound)))

	mux.HandleFunc("/notes", notesFunc)
	mux.HandleFunc("/notes/load-more", moreNotesFunc)

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
	mux.Handle("/favicon.ico", http.StripPrefix("/", http.FileServer(http.FS(fileSystem))))

	fmt.Println("Starting web interface on port: 8089")

	// Create stack for handle multiple middlewares
	middlewares := CreateMiddlewareStack(
		LoggingMiddleware,
		DemoMiddleware,
	)

	server := &http.Server{
		Addr:         ":8089",
		Handler:      middlewares(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
