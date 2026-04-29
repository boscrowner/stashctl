package snippet

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Mention represents a reference from one snippet to another, with an
// optional context note explaining why the link exists.
type Mention struct {
	ID          string    `json:"id"`
	FromID      string    `json:"from_id"`
	ToID        string    `json:"to_id"`
	Context     string    `json:"context,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

const maxMentionContext = 280

// NewMention creates a validated Mention linking fromID → toID.
func NewMention(fromID, toID, context string) (Mention, error) {
	fromID = strings.TrimSpace(fromID)
	toID = strings.TrimSpace(toID)
	context = strings.TrimSpace(context)

	if fromID == "" {
		return Mention{}, errors.New("mention: from_id must not be empty")
	}
	if toID == "" {
		return Mention{}, errors.New("mention: to_id must not be empty")
	}
	if fromID == toID {
		return Mention{}, errors.New("mention: a snippet cannot mention itself")
	}
	if len(context) > maxMentionContext {
		return Mention{}, errors.New("mention: context exceeds 280 characters")
	}

	return Mention{
		ID:        uuid.NewString(),
		FromID:    fromID,
		ToID:      toID,
		Context:   context,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// MentionsFrom returns all mentions whose FromID matches snippetID.
func MentionsFrom(mentions []Mention, snippetID string) []Mention {
	var out []Mention
	for _, m := range mentions {
		if m.FromID == snippetID {
			out = append(out, m)
		}
	}
	return out
}

// MentionsTo returns all mentions whose ToID matches snippetID.
func MentionsTo(mentions []Mention, snippetID string) []Mention {
	var out []Mention
	for _, m := range mentions {
		if m.ToID == snippetID {
			out = append(out, m)
		}
	}
	return out
}

// RemoveMention removes the mention with the given id from the slice.
func RemoveMention(mentions []Mention, id string) ([]Mention, error) {
	for i, m := range mentions {
		if m.ID == id {
			return append(mentions[:i], mentions[i+1:]...), nil
		}
	}
	return mentions, errors.New("mention: not found")
}
