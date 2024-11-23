package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sonjek/go-templ-htmx-picocss-example/internal/notes"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/templ/components"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/templ/page"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/templ/view"
)

func (h *Handlers) Notes(w http.ResponseWriter, r *http.Request) {
	handleRenderError(page.Index(view.NotesView(notes.GetLatestNotes())).Render(r.Context(), w))
}

func (h *Handlers) MoreNotes(w http.ResponseWriter, r *http.Request) {
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

func (h *Handlers) AddNoteModal(w http.ResponseWriter, r *http.Request) {
	handleRenderError(components.ModalAddNote().Render(r.Context(), w))
}

func (h *Handlers) AddNote(w http.ResponseWriter, r *http.Request) {
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

func (h *Handlers) EditNoteModal(w http.ResponseWriter, r *http.Request) {
	note, err := notes.GetNoteByID(r.PathValue("id"))
	if err != nil {
		sendErrorMsg(w, r, err.Error())
		return
	}

	handleRenderError(components.ModalEditNote(note).Render(r.Context(), w))
}

func (h *Handlers) EditNote(w http.ResponseWriter, r *http.Request) {
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

func (h *Handlers) DeleteNote(w http.ResponseWriter, r *http.Request) {
	if err := notes.Delete(r.PathValue("id")); err != nil {
		sendErrorMsg(w, r, err.Error())
	}
}
