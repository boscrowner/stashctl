package snippet

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TagGroup represents a named, ordered collection of tags that can be
// applied together to snippets.
type TagGroup struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
}

// NewTagGroup creates a validated TagGroup.
func NewTagGroup(name, description string, tags []string) (TagGroup, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return TagGroup{}, errors.New("tag group name must not be empty")
	}
	if len(name) > 64 {
		return TagGroup{}, errors.New("tag group name must not exceed 64 characters")
	}
	normalized := NormalizeTags(tags)
	if len(normalized) == 0 {
		return TagGroup{}, errors.New("tag group must contain at least one tag")
	}
	return TagGroup{
		ID:          uuid.NewString(),
		Name:        name,
		Description: strings.TrimSpace(description),
		Tags:        normalized,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

// TagGroupsFor returns all tag groups from the slice.
func TagGroupsFor(groups []TagGroup) []TagGroup {
	out := make([]TagGroup, len(groups))
	copy(out, groups)
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

// FindTagGroup returns the first tag group matching the given name (case-insensitive).
func FindTagGroup(groups []TagGroup, name string) (TagGroup, bool) {
	norm := strings.ToLower(strings.TrimSpace(name))
	for _, g := range groups {
		if strings.ToLower(g.Name) == norm {
			return g, true
		}
	}
	return TagGroup{}, false
}

// RemoveTagGroup returns a new slice without the group with the given id.
func RemoveTagGroup(groups []TagGroup, id string) []TagGroup {
	out := groups[:0:0]
	for _, g := range groups {
		if g.ID != id {
			out = append(out, g)
		}
	}
	return out
}
