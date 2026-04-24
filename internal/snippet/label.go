package snippet

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Label represents a named colour label that can be attached to snippets.
type Label struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"` // hex, e.g. "#ff5733"
	CreatedAt time.Time `json:"created_at"`
}

// NewLabel creates a new Label with the given name and colour.
func NewLabel(name, color string) (Label, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Label{}, errors.New("label name must not be empty")
	}
	if len(name) > 32 {
		return Label{}, errors.New("label name must not exceed 32 characters")
	}
	color = strings.TrimSpace(color)
	if color == "" {
		color = "#cccccc"
	}
	return Label{
		ID:        uuid.NewString(),
		Name:      name,
		Color:     color,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// LabelsFor returns all labels whose IDs are present in ids.
func LabelsFor(all []Label, ids []string) []Label {
	set := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	var out []Label
	for _, l := range all {
		if _, ok := set[l.ID]; ok {
			out = append(out, l)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

// RemoveLabel removes the label with the given id from the slice.
func RemoveLabel(all []Label, id string) []Label {
	out := all[:0:0]
	for _, l := range all {
		if l.ID != id {
			out = append(out, l)
		}
	}
	return out
}

// FindLabel returns the first label matching name (case-insensitive) and true,
// or an empty Label and false if not found.
func FindLabel(all []Label, name string) (Label, bool) {
	name = strings.ToLower(strings.TrimSpace(name))
	for _, l := range all {
		if strings.ToLower(l.Name) == name {
			return l, true
		}
	}
	return Label{}, false
}
