package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/templ/components"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/templ/page"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/templ/view"
)

const (
	pageSize  = 2
	timeoutMs = 300
)

func (h *Handlers) Notes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.noteService.LoadMore(0, pageSize)
	if err != nil {
		sendErrorMsg(w, r, "Note is empty")
	}

	// Timeout for show loader
	time.Sleep(timeoutMs * time.Millisecond)

	handleRenderError(page.Index(view.NotesView(notes)).Render(r.Context(), w))
}

func (h *Handlers) LoadMoreNotes(w http.ResponseWriter, r *http.Request) {
	cursor := -1
	if p := r.URL.Query().Get("cursor"); p != "" {
		if parsedCursor, err := strconv.Atoi(p); err == nil {
			cursor = parsedCursor
		}
	}

	notes, err := h.noteService.LoadMore(cursor, pageSize)
	if err != nil {
		sendErrorMsg(w, r, "Note is empty")
	}

	// Timeout for show loader
	time.Sleep(timeoutMs * time.Millisecond)

	handleRenderError(components.NotesList(notes).Render(r.Context(), w))
}

func (h *Handlers) CreateNoteModal(w http.ResponseWriter, r *http.Request) {
	handleRenderError(components.ModalAddNote().Render(r.Context(), w))
}

func (h *Handlers) CreateNote(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	if title == "" {
		sendErrorMsg(w, r, "Title is empty")
		return
	}

	body := r.FormValue("body")
	if body == "" {
		sendErrorMsg(w, r, "Body is empty")
		return
	}

	note := h.noteService.Create(title, body)

	// Timeout for show loader
	time.Sleep(timeoutMs * time.Millisecond)

	handleRenderError(components.NoteItem(note).Render(r.Context(), w))
}

func (h *Handlers) EditNoteModal(w http.ResponseWriter, r *http.Request) {
	noteID := -1
	if p := r.PathValue("id"); p != "" {
		if parsedNoteID, err := strconv.Atoi(p); err == nil {
			noteID = parsedNoteID
		}
	}

	if noteID < 1 {
		sendErrorMsg(w, r, "Wrong note ID")
		return
	}

	note := h.noteService.Get(noteID)

	handleRenderError(components.ModalEditNote(note).Render(r.Context(), w))
}

func (h *Handlers) EditNote(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	if title == "" {
		sendErrorMsg(w, r, "Title is empty")
		return
	}

	body := r.FormValue("body")
	if body == "" {
		sendErrorMsg(w, r, "Body is empty")
		return
	}

	noteID := -1
	if p := r.PathValue("id"); p != "" {
		if parsedNoteID, err := strconv.Atoi(p); err == nil {
			noteID = parsedNoteID
		}
	}

	note := h.noteService.FindAndUpdate(noteID, title, body)

	// Timeout for show loader
	time.Sleep(timeoutMs * time.Millisecond)

	handleRenderError(components.NoteItem(note).Render(r.Context(), w))
}

func (h *Handlers) DeleteNote(w http.ResponseWriter, r *http.Request) {
	noteID := r.PathValue("id")
	if noteID == "" {
		sendErrorMsg(w, r, "Note ID is empty")
		return
	}

	// Timeout for show loader
	time.Sleep(timeoutMs * time.Millisecond)

	h.noteService.Delete(noteID)
}
