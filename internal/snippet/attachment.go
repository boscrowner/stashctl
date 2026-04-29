package snippet

import (
	"errors"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Attachment represents a file reference associated with a snippet.
type Attachment struct {
	ID        string    `json:"id"`
	SnippetID string    `json:"snippet_id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	MIMEType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
}

// NewAttachment creates a validated Attachment for the given snippet.
func NewAttachment(snippetID, name, path, mimeType string) (Attachment, error) {
	snippetID = strings.TrimSpace(snippetID)
	name = strings.TrimSpace(name)
	path = strings.TrimSpace(path)

	if snippetID == "" {
		return Attachment{}, errors.New("attachment: snippet_id must not be empty")
	}
	if name == "" {
		return Attachment{}, errors.New("attachment: name must not be empty")
	}
	if len(name) > 128 {
		return Attachment{}, errors.New("attachment: name must not exceed 128 characters")
	}
	if path == "" {
		return Attachment{}, errors.New("attachment: path must not be empty")
	}
	if !filepath.IsAbs(path) && !strings.HasPrefix(path, ".") {
		path = "." + string(filepath.Separator) + path
	}

	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return Attachment{
		ID:        uuid.NewString(),
		SnippetID: snippetID,
		Name:      name,
		Path:      path,
		MIMEType:  mimeType,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// AttachmentsFor returns all attachments belonging to the given snippet ID.
func AttachmentsFor(attachments []Attachment, snippetID string) []Attachment {
	var out []Attachment
	for _, a := range attachments {
		if a.SnippetID == snippetID {
			out = append(out, a)
		}
	}
	return out
}

// RemoveAttachment removes the attachment with the given ID from the slice.
// It returns the updated slice and true if an attachment was removed.
func RemoveAttachment(attachments []Attachment, id string) ([]Attachment, bool) {
	for i, a := range attachments {
		if a.ID == id {
			return append(attachments[:i], attachments[i+1:]...), true
		}
	}
	return attachments, false
}

// FindAttachment returns the attachment with the given ID and true if found,
// or a zero-value Attachment and false if no match exists.
func FindAttachment(attachments []Attachment, id string) (Attachment, bool) {
	for _, a := range attachments {
		if a.ID == id {
			return a, true
		}
	}
	return Attachment{}, false
}
