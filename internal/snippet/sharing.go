package snippet

import (
	"errors"
	"fmt"
	"time"

	"github.com/user/stashctl/internal/snippet"
)

// ShareVisibility controls who can access a shared snippet.
type ShareVisibility string

const (
	VisibilityPublic  ShareVisibility = "public"
	VisibilityPrivate ShareVisibility = "private"
	VisibilityUnlisted ShareVisibility = "unlisted"
)

// ShareLink represents a shareable link for a snippet.
type ShareLink struct {
	ID        string          `json:"id"`
	SnippetID string          `json:"snippet_id"`
	Token     string          `json:"token"`
	Visibility ShareVisibility `json:"visibility"`
	ExpiresAt *time.Time      `json:"expires_at,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	Note      string          `json:"note,omitempty"`
}

// NewShareLink creates a new ShareLink for the given snippet ID.
func NewShareLink(snippetID string, visibility ShareVisibility, expiresAt *time.Time, note string) (ShareLink, error) {
	if snippetID == "" {
		return ShareLink{}, errors.New("snippet ID must not be empty")
	}
	if visibility == "" {
		visibility = VisibilityPrivate
	}
	if !isKnownVisibility(visibility) {
		return ShareLink{}, fmt.Errorf("unknown visibility %q: must be public, private, or unlisted", visibility)
	}
	now := time.Now().UTC()
	if expiresAt != nil && expiresAt.Before(now) {
		return ShareLink{}, errors.New("expiry time must be in the future")
	}
	if len(note) > 200 {
		return ShareLink{}, errors.New("note must not exceed 200 characters")
	}
	return ShareLink{
		ID:         generateID(),
		SnippetID:  snippetID,
		Token:      generateID(),
		Visibility: visibility,
		ExpiresAt:  expiresAt,
		CreatedAt:  now,
		Note:       note,
	}, nil
}

// IsExpired reports whether the share link has passed its expiry time.
func (s ShareLink) IsExpired() bool {
	if s.ExpiresAt == nil {
		return false
	}
	return time.Now().UTC().After(*s.ExpiresAt)
}

// ShareLinksFor returns all share links associated with the given snippet ID.
func ShareLinksFor(links []ShareLink, snippetID string) []ShareLink {
	var result []ShareLink
	for _, l := range links {
		if l.SnippetID == snippetID {
			result = append(result, l)
		}
	}
	return result
}

// RemoveShareLink removes the share link with the given ID from the slice.
func RemoveShareLink(links []ShareLink, id string) ([]ShareLink, error) {
	for i, l := range links {
		if l.ID == id {
			return append(links[:i], links[i+1:]...), nil
		}
	}
	return links, fmt.Errorf("share link %q not found", id)
}

func isKnownVisibility(v ShareVisibility) bool {
	switch v {
	case VisibilityPublic, VisibilityPrivate, VisibilityUnlisted:
		return true
	}
	return false
}
