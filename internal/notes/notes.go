package notes

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/dustin/go-humanize"
)

type Note struct {
	ID      int       `form:"id" json:"id"`
	Title   string    `form:"title" json:"title" binding:"required"`
	Body    string    `form:"body" json:"body" binding:"required"`
	Created time.Time `json:"created"`
}

type CreateNote struct {
	Title string `form:"title" json:"title" binding:"required"`
	Body  string `form:"body" json:"body" binding:"required"`
}

var (
	notes = []Note{
		{
			ID:      5,
			Title:   "htmx.js",
			Body:    "</> htmx - high power tools for HTML",
			Created: time.Date(2022, 8, 10, 21, 21, 0, 0, time.UTC),
		},
		{
			ID:      4,
			Title:   "templ",
			Body:    "A language for writing HTML user interfaces in Go.",
			Created: time.Date(2021, 5, 16, 22, 33, 0, 0, time.UTC),
		},
		{
			ID:      3,
			Title:   "Pico CSS",
			Body:    "Minimal CSS Framework for semantic HTML",
			Created: time.Date(2019, 12, 11, 10, 8, 0, 0, time.UTC),
		},
		{
			ID:      2,
			Title:   "GoLang",
			Body:    "The Go programming language.",
			Created: time.Date(2018, 9, 25, 22, 20, 0, 0, time.UTC),
		},
		{
			ID:      1,
			Title:   "Ionic",
			Body:    "Premium hand-crafted icons built by Ionic, for Ionic apps and web apps everywhere.",
			Created: time.Date(2013, 10, 30, 12, 34, 0, 0, time.UTC),
		},
	}

	// Last note ID
	currentID = uint32(len(notes))
)

func getNextID() uint32 {
	return atomic.AddUint32(&currentID, 1)
}

func GetNoteByID(idStr string) (Note, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return Note{}, fmt.Errorf("wrong ID: %s", idStr)
	}

	for _, note := range notes {
		if note.ID == id {
			return note, nil
		}
	}
	return Note{}, fmt.Errorf("note with ID %d not found", id)
}

func (n Note) FormatCreatedAgo() string {
	return humanize.Time(n.Created)
}

func (n Note) FormatCreated() string {
	return n.Created.Format("2006-01-02 15:04")
}

func GetAll() []Note {
	return notes
}

func Count() int {
	return len(notes)
}

func Add(n CreateNote) Note {
	note := Note{
		ID:      int(getNextID()),
		Title:   n.Title,
		Body:    n.Body,
		Created: time.Now(),
	}
	notes = append([]Note{note}, notes...)
	return note
}

func Update(n Note) {
	for i, note := range notes {
		if note.ID == n.ID {
			notes[i] = n
			break
		}
	}
}

func Delete(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("wrong ID: %s", idStr)
	}

	for idx, note := range notes {
		if note.ID == id {
			notes = append(notes[:idx], notes[idx+1:]...)
			return nil
		}
	}
	return fmt.Errorf("note with ID %d not found", id)
}

func findIndex(arr []Note, n int) (int, bool) {
	index := -1

	// Initialize to -1 to represent no ID found yet
	maxID := -1

	for i, elem := range arr {
		if elem.ID < n && elem.ID > maxID {
			index = i
			maxID = elem.ID
		}
	}

	return index, index != -1
}

func GetNextNotes(noteID int) []Note {
	startIndex, found := findIndex(notes, noteID)
	if !found {
		return []Note{}
	}

	available := notes[startIndex:]

	// Adjust endIndex if it exceeds the length of items
	endIndex := 3
	if endIndex > len(available) {
		endIndex = len(available)
	}
	return available[:endIndex]
}

func GetLatestNotes() []Note {
	return GetNextNotes(int(currentID + 1))
}
