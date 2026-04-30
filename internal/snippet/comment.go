package snippet

import (
	"errors"
	"time"
)

const maxCommentBody = 1000

// Comment represents a user comment attached to a snippet.
type Comment struct {
	ID        string    `json:"id"`
	SnippetID string    `json:"snippet_id"`
	Author    string    `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewComment creates a new Comment for the given snippet.
func NewComment(snippetID, author, body string) (Comment, error) {
	if snippetID == "" {
		return Comment{}, errors.New("snippet ID must not be empty")
	}
	if author == "" {
		return Comment{}, errors.New("author must not be empty")
	}
	if body == "" {
		return Comment{}, errors.New("body must not be empty")
	}
	if len(body) > maxCommentBody {
		return Comment{}, errors.New("body exceeds maximum length")
	}
	now := time.Now().UTC()
	return Comment{
		ID:        generateID(),
		SnippetID: snippetID,
		Author:    author,
		Body:      body,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// CommentsFor returns all comments belonging to the given snippet ID.
func CommentsFor(snippetID string, all []Comment) []Comment {
	var out []Comment
	for _, c := range all {
		if c.SnippetID == snippetID {
			out = append(out, c)
		}
	}
	return out
}

// RemoveComment removes a comment by ID, returning the updated slice and
// a boolean indicating whether the comment was found.
func RemoveComment(id string, all []Comment) ([]Comment, bool) {
	out := make([]Comment, 0, len(all))
	found := false
	for _, c := range all {
		if c.ID == id {
			found = true
			continue
		}
		out = append(out, c)
	}
	return out, found
}

// FindComment returns the comment with the given ID, if it exists.
func FindComment(id string, all []Comment) (Comment, bool) {
	for _, c := range all {
		if c.ID == id {
			return c, true
		}
	}
	return Comment{}, false
}
