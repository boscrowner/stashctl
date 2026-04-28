package snippet

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Alias maps a short human-readable name to a snippet ID.
type Alias struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	SnippetID string    `json:"snippet_id"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

const maxAliasNameLen = 64

// NewAlias creates a validated Alias linking name to snippetID.
func NewAlias(name, snippetID, note string) (Alias, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Alias{}, errors.New("alias name must not be empty")
	}
	if len(name) > maxAliasNameLen {
		return Alias{}, errors.New("alias name exceeds maximum length")
	}
	if strings.ContainsAny(name, " \t\n") {
		return Alias{}, errors.New("alias name must not contain whitespace")
	}
	if snippetID == "" {
		return Alias{}, errors.New("snippet ID must not be empty")
	}
	if len(note) > 256 {
		return Alias{}, errors.New("note exceeds maximum length of 256 characters")
	}
	return Alias{
		ID:        uuid.NewString(),
		Name:      name,
		SnippetID: snippetID,
		Note:      strings.TrimSpace(note),
		CreatedAt: time.Now().UTC(),
	}, nil
}

// AliasesFor returns all aliases associated with the given snippet ID.
func AliasesFor(aliases []Alias, snippetID string) []Alias {
	var out []Alias
	for _, a := range aliases {
		if a.SnippetID == snippetID {
			out = append(out, a)
		}
	}
	return out
}

// FindAlias returns the alias with the given name, or false if not found.
func FindAlias(aliases []Alias, name string) (Alias, bool) {
	for _, a := range aliases {
		if a.Name == name {
			return a, true
		}
	}
	return Alias{}, false
}

// RemoveAlias returns a new slice with the alias identified by id removed.
func RemoveAlias(aliases []Alias, id string) ([]Alias, error) {
	out := make([]Alias, 0, len(aliases))
	found := false
	for _, a := range aliases {
		if a.ID == id {
			found = true
			continue
		}
		out = append(out, a)
	}
	if !found {
		return nil, errors.New("alias not found")
	}
	return out, nil
}
