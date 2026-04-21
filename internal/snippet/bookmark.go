package snippet

import (
	"fmt"
	"time"
)

const bookmarkTag = "__bookmark__"

// Bookmark represents a named reference to a snippet with optional notes.
type Bookmark struct {
	SnippetID string    `json:"snippet_id"`
	Name      string    `json:"name"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// NewBookmark creates a new Bookmark for the given snippet ID and name.
// Returns an error if snippetID or name is empty.
func NewBookmark(snippetID, name, note string) (Bookmark, error) {
	if snippetID == "" {
		return Bookmark{}, fmt.Errorf("snippet ID must not be empty")
	}
	if name == "" {
		return Bookmark{}, fmt.Errorf("bookmark name must not be empty")
	}
	if len(note) > 200 {
		return Bookmark{}, fmt.Errorf("note must not exceed 200 characters")
	}
	return Bookmark{
		SnippetID: snippetID,
		Name:      name,
		Note:      note,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// BookmarksFor returns all bookmarks associated with the given snippet ID.
func BookmarksFor(snippetID string, bookmarks []Bookmark) []Bookmark {
	var result []Bookmark
	for _, b := range bookmarks {
		if b.SnippetID == snippetID {
			result = append(result, b)
		}
	}
	return result
}

// RemoveBookmark removes the bookmark with the given name for a snippet.
// Returns the updated slice and true if a bookmark was removed.
func RemoveBookmark(name, snippetID string, bookmarks []Bookmark) ([]Bookmark, bool) {
	updated := make([]Bookmark, 0, len(bookmarks))
	removed := false
	for _, b := range bookmarks {
		if b.Name == name && b.SnippetID == snippetID {
			removed = true
			continue
		}
		updated = append(updated, b)
	}
	return updated, removed
}

// FindBookmark returns the first bookmark matching the given name across all snippets.
func FindBookmark(name string, bookmarks []Bookmark) (Bookmark, bool) {
	for _, b := range bookmarks {
		if b.Name == name {
			return b, true
		}
	}
	return Bookmark{}, false
}
