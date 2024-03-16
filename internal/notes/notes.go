package notes

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
)

type Note struct {
	Id      int       `form:"id" json:"id"`
	Title   string    `form:"title" json:"title" binding:"required"`
	Body    string    `form:"body" json:"body" binding:"required"`
	Created time.Time `json:"created"`
}

type CreateNote struct {
	Title string `form:"title" json:"title" binding:"required"`
	Body  string `form:"body" json:"body" binding:"required"`
}

var notes []Note = []Note{
	{
		Id:      4,
		Title:   "htmx.js",
		Body:    "</> htmx - high power tools for HTML",
		Created: time.Date(2022, 8, 10, 21, 21, 0, 0, time.UTC),
	},
	{
		Id:      3,
		Title:   "templ",
		Body:    "A language for writing HTML user interfaces in Go.",
		Created: time.Date(2021, 5, 16, 22, 33, 0, 0, time.UTC),
	},
	{
		Id:      2,
		Title:   "Pico CSS",
		Body:    "Minimal CSS Framework for semantic HTML",
		Created: time.Date(2019, 12, 11, 10, 8, 0, 0, time.UTC),
	},
	{
		Id:      1,
		Title:   "GoLang",
		Body:    "The Go programming language.",
		Created: time.Date(2018, 9, 25, 22, 20, 0, 0, time.UTC),
	},
}

// Last note ID
var currentID uint32 = 4

func getNextID() uint32 {
	return atomic.AddUint32(&currentID, 1)
}

func GetNoteByID(idStr string) (Note, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return Note{}, fmt.Errorf("wrong ID: %s", idStr)
	}

	for _, note := range notes {
		if note.Id == id {
			return note, nil
		}
	}
	return Note{}, fmt.Errorf("note with id %d not found", id)
}

func GetAll() []Note {
	return notes
}

func Count() int {
	return len(notes)
}

func Add(n CreateNote) {
	note := Note{
		Id:      int(getNextID()),
		Title:   n.Title,
		Body:    n.Body,
		Created: time.Now(),
	}
	notes = append([]Note{note}, notes...)
}

func Update(n Note) {
	for i, note := range notes {
		if note.Id == n.Id {
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
		if note.Id == id {
			notes = append(notes[:idx], notes[idx+1:]...)
			return nil
		}

	}
	return fmt.Errorf("note with id %d not found", id)
}
