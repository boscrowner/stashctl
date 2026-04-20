package snippet

import (
	"errors"
	"time"
)

// Collection groups snippets under a named label with optional description.
type Collection struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	SnippetIDs  []string  `json:"snippet_ids"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewCollection creates a new Collection with the given name and description.
func NewCollection(name, description string) (*Collection, error) {
	if name == "" {
		return nil, errors.New("collection name must not be empty")
	}
	if len(name) > 100 {
		return nil, errors.New("collection name must not exceed 100 characters")
	}
	now := time.Now().UTC()
	return &Collection{
		ID:          generateID(),
		Name:        name,
		Description: description,
		SnippetIDs:  []string{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// AddSnippet adds a snippet ID to the collection if not already present.
func (c *Collection) AddSnippet(id string) {
	for _, existing := range c.SnippetIDs {
		if existing == id {
			return
		}
	}
	c.SnippetIDs = append(c.SnippetIDs, id)
	c.UpdatedAt = time.Now().UTC()
}

// RemoveSnippet removes a snippet ID from the collection.
// Returns true if the snippet was found and removed.
func (c *Collection) RemoveSnippet(id string) bool {
	for i, existing := range c.SnippetIDs {
		if existing == id {
			c.SnippetIDs = append(c.SnippetIDs[:i], c.SnippetIDs[i+1:]...)
			c.UpdatedAt = time.Now().UTC()
			return true
		}
	}
	return false
}

// Contains reports whether the collection holds the given snippet ID.
func (c *Collection) Contains(id string) bool {
	for _, existing := range c.SnippetIDs {
		if existing == id {
			return true
		}
	}
	return false
}
