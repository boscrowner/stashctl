package snippet

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Workspace groups snippets into a named working context.
type Workspace struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SnippetIDs  []string  `json:"snippet_ids"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewWorkspace creates a new Workspace with the given name and description.
func NewWorkspace(name, description string) (*Workspace, error) {
	if name == "" {
		return nil, errors.New("workspace name must not be empty")
	}
	if len(name) > 64 {
		return nil, errors.New("workspace name must not exceed 64 characters")
	}
	now := time.Now().UTC()
	return &Workspace{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		SnippetIDs:  []string{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// AddSnippet adds a snippet ID to the workspace if not already present.
func (w *Workspace) AddSnippet(snippetID string) error {
	if snippetID == "" {
		return errors.New("snippet ID must not be empty")
	}
	for _, id := range w.SnippetIDs {
		if id == snippetID {
			return nil
		}
	}
	w.SnippetIDs = append(w.SnippetIDs, snippetID)
	w.UpdatedAt = time.Now().UTC()
	return nil
}

// RemoveSnippet removes a snippet ID from the workspace.
func (w *Workspace) RemoveSnippet(snippetID string) error {
	for i, id := range w.SnippetIDs {
		if id == snippetID {
			w.SnippetIDs = append(w.SnippetIDs[:i], w.SnippetIDs[i+1:]...)
			w.UpdatedAt = time.Now().UTC()
			return nil
		}
	}
	return errors.New("snippet not found in workspace")
}

// WorkspacesFor returns all workspaces that contain the given snippet ID.
func WorkspacesFor(snippetID string, workspaces []*Workspace) []*Workspace {
	var result []*Workspace
	for _, w := range workspaces {
		for _, id := range w.SnippetIDs {
			if id == snippetID {
				result = append(result, w)
				break
			}
		}
	}
	return result
}
