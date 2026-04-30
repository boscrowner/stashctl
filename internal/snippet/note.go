package snippet

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Note is a free-form text note attached to a snippet.
type Note struct {
	ID        string    `json:"id"`
	SnippetID string    `json:"snippet_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const maxNoteBody = 2000

// NewNote creates a new Note attached to snippetID.
func NewNote(snippetID, body string) (Note, error) {
	if snippetID == "" {
		return Note{}, errors.New("note: snippet ID must not be empty")
	}
	if body == "" {
		return Note{}, errors.New("note: body must not be empty")
	}
	if len(body) > maxNoteBody {
		return Note{}, errors.New("note: body exceeds 2000 characters")
	}
	now := time.Now().UTC()
	return Note{
		ID:        uuid.NewString(),
		SnippetID: snippetID,
		Body:      body,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// NotesFor returns all notes belonging to snippetID.
func NotesFor(notes []Note, snippetID string) []Note {
	var out []Note
	for _, n := range notes {
		if n.SnippetID == snippetID {
			out = append(out, n)
		}
	}
	return out
}

// RemoveNote removes the note with the given id from the slice.
// Returns the updated slice and true if a note was removed.
func RemoveNote(notes []Note, id string) ([]Note, bool) {
	for i, n := range notes {
		if n.ID == id {
			return append(notes[:i], notes[i+1:]...), true
		}
	}
	return notes, false
}

// FindNote returns the note with the given id, if present.
func FindNote(notes []Note, id string) (Note, bool) {
	for _, n := range notes {
		if n.ID == id {
			return n, true
		}
	}
	return Note{}, false
}
