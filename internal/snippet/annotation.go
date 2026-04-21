package snippet

import (
	"errors"
	"strings"
	"time"
)

// Annotation represents a note attached to a snippet.
type Annotation struct {
	ID        string    `json:"id"`
	SnippetID string    `json:"snippet_id"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

// NewAnnotation creates a new Annotation for the given snippet ID and note text.
func NewAnnotation(snippetID, note string) (Annotation, error) {
	snippetID = strings.TrimSpace(snippetID)
	if snippetID == "" {
		return Annotation{}, errors.New("snippet ID must not be empty")
	}
	note = strings.TrimSpace(note)
	if note == "" {
		return Annotation{}, errors.New("note must not be empty")
	}
	if len(note) > 500 {
		return Annotation{}, errors.New("note must not exceed 500 characters")
	}
	return Annotation{
		ID:        generateID(),
		SnippetID: snippetID,
		Note:      note,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// AnnotationsFor returns all annotations belonging to the given snippet ID.
func AnnotationsFor(annotations []Annotation, snippetID string) []Annotation {
	var result []Annotation
	for _, a := range annotations {
		if a.SnippetID == snippetID {
			result = append(result, a)
		}
	}
	return result
}

// RemoveAnnotation removes the annotation with the given ID from the slice.
// Returns the updated slice and true if an annotation was removed.
func RemoveAnnotation(annotations []Annotation, id string) ([]Annotation, bool) {
	for i, a := range annotations {
		if a.ID == id {
			return append(annotations[:i], annotations[i+1:]...), true
		}
	}
	return annotations, false
}
