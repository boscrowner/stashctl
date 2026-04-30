package snippet

import (
	"fmt"
	"strings"
	"time"
)

// Badge represents a visual achievement or status marker attached to a snippet.
type Badge struct {
	ID        string    `json:"id"`
	SnippetID string    `json:"snippet_id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	Color     string    `json:"color"`
	GrantedAt time.Time `json:"granted_at"`
	GrantedBy string    `json:"granted_by"`
	Note      string    `json:"note"`
}

const (
	maxBadgeNameLen = 40
	maxBadgeNoteLen = 200
)

var knownBadgeColors = map[string]bool{
	"red": true, "green": true, "blue": true,
	"yellow": true, "purple": true, "orange": true,
	"gray": true, "gold": true,
}

// NewBadge creates a validated Badge for the given snippet.
func NewBadge(snippetID, name, icon, color, grantedBy, note string) (Badge, error) {
	snippetID = strings.TrimSpace(snippetID)
	name = strings.TrimSpace(name)
	grantedBy = strings.TrimSpace(grantedBy)

	if snippetID == "" {
		return Badge{}, fmt.Errorf("badge: snippet ID must not be empty")
	}
	if name == "" {
		return Badge{}, fmt.Errorf("badge: name must not be empty")
	}
	if len(name) > maxBadgeNameLen {
		return Badge{}, fmt.Errorf("badge: name exceeds %d characters", maxBadgeNameLen)
	}
	if grantedBy == "" {
		return Badge{}, fmt.Errorf("badge: granted_by must not be empty")
	}
	if len(note) > maxBadgeNoteLen {
		return Badge{}, fmt.Errorf("badge: note exceeds %d characters", maxBadgeNoteLen)
	}
	if color != "" && !knownBadgeColors[color] {
		return Badge{}, fmt.Errorf("badge: unknown color %q", color)
	}
	if color == "" {
		color = "gray"
	}

	return Badge{
		ID:        generateID(),
		SnippetID: snippetID,
		Name:      name,
		Icon:      icon,
		Color:     color,
		GrantedAt: time.Now().UTC(),
		GrantedBy: grantedBy,
		Note:      note,
	}, nil
}

// BadgesFor returns all badges associated with the given snippet ID.
func BadgesFor(snippetID string, badges []Badge) []Badge {
	var out []Badge
	for _, b := range badges {
		if b.SnippetID == snippetID {
			out = append(out, b)
		}
	}
	return out
}

// RemoveBadge returns a new slice with the badge matching id removed.
func RemoveBadge(id string, badges []Badge) ([]Badge, bool) {
	out := make([]Badge, 0, len(badges))
	removed := false
	for _, b := range badges {
		if b.ID == id {
			removed = true
			continue
		}
		out = append(out, b)
	}
	return out, removed
}

// FindBadge returns the badge with the given id, if present.
func FindBadge(id string, badges []Badge) (Badge, bool) {
	for _, b := range badges {
		if b.ID == id {
			return b, true
		}
	}
	return Badge{}, false
}
