package snippet

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// KnownReactions is the set of allowed reaction emoji slugs.
var KnownReactions = map[string]bool{
	"thumbsup":   true,
	"thumbsdown": true,
	"heart":      true,
	"fire":       true,
	"rocket":     true,
	"eyes":       true,
	"tada":       true,
	"confused":   true,
}

// Reaction represents a single emoji reaction left by a user on a snippet.
type Reaction struct {
	ID        string    `json:"id"`
	SnippetID string    `json:"snippet_id"`
	Actor     string    `json:"actor"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
}

// NewReaction creates and validates a new Reaction.
func NewReaction(snippetID, actor, emoji string) (Reaction, error) {
	if snippetID == "" {
		return Reaction{}, errors.New("reaction: snippet_id is required")
	}
	if actor == "" {
		return Reaction{}, errors.New("reaction: actor is required")
	}
	if !KnownReactions[emoji] {
		return Reaction{}, errors.New("reaction: unknown emoji slug")
	}
	return Reaction{
		ID:        uuid.NewString(),
		SnippetID: snippetID,
		Actor:     actor,
		Emoji:     emoji,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// ReactionsFor returns all reactions for the given snippet ID.
func ReactionsFor(reactions []Reaction, snippetID string) []Reaction {
	var out []Reaction
	for _, r := range reactions {
		if r.SnippetID == snippetID {
			out = append(out, r)
		}
	}
	return out
}

// CountByEmoji returns a map of emoji slug → count for the given reactions.
func CountByEmoji(reactions []Reaction) map[string]int {
	counts := make(map[string]int)
	for _, r := range reactions {
		counts[r.Emoji]++
	}
	return counts
}

// RemoveReaction removes the reaction with the given ID from the slice.
func RemoveReaction(reactions []Reaction, id string) ([]Reaction, error) {
	for i, r := range reactions {
		if r.ID == id {
			return append(reactions[:i], reactions[i+1:]...), nil
		}
	}
	return reactions, errors.New("reaction: not found")
}
