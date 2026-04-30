package snippet

import (
	"errors"
	"net/url"
	"time"

	"github.com/google/uuid"
)

// Link represents a URL associated with a snippet.
type Link struct {
	ID        string    `json:"id"`
	SnippetID string    `json:"snippet_id"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

// NewLink creates a new Link for the given snippet.
// url must be a valid HTTP/HTTPS URL. title is optional but must not exceed 120 chars.
func NewLink(snippetID, rawURL, title string) (Link, error) {
	if snippetID == "" {
		return Link{}, errors.New("snippet id must not be empty")
	}
	if rawURL == "" {
		return Link{}, errors.New("url must not be empty")
	}
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return Link{}, errors.New("url is not valid")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return Link{}, errors.New("url scheme must be http or https")
	}
	if len(title) > 120 {
		return Link{}, errors.New("title must not exceed 120 characters")
	}
	return Link{
		ID:        uuid.NewString(),
		SnippetID: snippetID,
		URL:       rawURL,
		Title:     title,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// LinksFor returns all links associated with the given snippet ID.
func LinksFor(links []Link, snippetID string) []Link {
	var out []Link
	for _, l := range links {
		if l.SnippetID == snippetID {
			out = append(out, l)
		}
	}
	return out
}

// RemoveLink removes the link with the given id from the slice.
// Returns the updated slice and true if a link was removed.
func RemoveLink(links []Link, id string) ([]Link, bool) {
	for i, l := range links {
		if l.ID == id {
			return append(links[:i], links[i+1:]...), true
		}
	}
	return links, false
}

// FindLink returns the link with the given id, or false if not found.
func FindLink(links []Link, id string) (Link, bool) {
	for _, l := range links {
		if l.ID == id {
			return l, true
		}
	}
	return Link{}, false
}
