package snippet

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// AccessLevel represents the permission level for a snippet.
type AccessLevel string

const (
	AccessRead  AccessLevel = "read"
	AccessWrite AccessLevel = "write"
	AccessAdmin AccessLevel = "admin"
)

// AccessGrant records that a principal has been granted access to a snippet.
type AccessGrant struct {
	ID        string      `json:"id"`
	SnippetID string      `json:"snippet_id"`
	Principal string      `json:"principal"`
	Level     AccessLevel `json:"level"`
	GrantedAt time.Time   `json:"granted_at"`
	Note      string      `json:"note,omitempty"`
}

var knownAccessLevels = map[AccessLevel]bool{
	AccessRead:  true,
	AccessWrite: true,
	AccessAdmin: true,
}

// NewAccessGrant creates a validated AccessGrant for the given snippet and principal.
func NewAccessGrant(snippetID, principal string, level AccessLevel, note string) (AccessGrant, error) {
	if snippetID == "" {
		return AccessGrant{}, errors.New("snippet ID must not be empty")
	}
	if principal == "" {
		return AccessGrant{}, errors.New("principal must not be empty")
	}
	if !knownAccessLevels[level] {
		return AccessGrant{}, errors.New("unknown access level: " + string(level))
	}
	if len(note) > 200 {
		return AccessGrant{}, errors.New("note must not exceed 200 characters")
	}
	return AccessGrant{
		ID:        uuid.NewString(),
		SnippetID: snippetID,
		Principal: principal,
		Level:     level,
		GrantedAt: time.Now().UTC(),
		Note:      note,
	}, nil
}

// AccessGrantsFor returns all grants associated with the given snippet ID.
func AccessGrantsFor(grants []AccessGrant, snippetID string) []AccessGrant {
	var out []AccessGrant
	for _, g := range grants {
		if g.SnippetID == snippetID {
			out = append(out, g)
		}
	}
	return out
}

// RemoveAccessGrant removes the grant with the given ID from the slice.
func RemoveAccessGrant(grants []AccessGrant, id string) ([]AccessGrant, error) {
	for i, g := range grants {
		if g.ID == id {
			return append(grants[:i], grants[i+1:]...), nil
		}
	}
	return grants, errors.New("access grant not found: " + id)
}

// HasAccess reports whether any grant for snippetID covers principal at the requested level.
func HasAccess(grants []AccessGrant, snippetID, principal string, required AccessLevel) bool {
	order := map[AccessLevel]int{AccessRead: 1, AccessWrite: 2, AccessAdmin: 3}
	for _, g := range grants {
		if g.SnippetID == snippetID && g.Principal == principal {
			if order[g.Level] >= order[required] {
				return true
			}
		}
	}
	return false
}
