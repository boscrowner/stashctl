package snippet

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Snapshot captures the full state of a snippet at a point in time.
type Snapshot struct {
	ID        string    `json:"id"`
	SnippetID string    `json:"snippet_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Language  string    `json:"language"`
	Tags      []string  `json:"tags"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"created_at"`
}

// NewSnapshot creates a snapshot of the given snippet.
func NewSnapshot(s Snippet, label string) (Snapshot, error) {
	if s.ID == "" {
		return Snapshot{}, errors.New("snippet id must not be empty")
	}
	if len(label) > 64 {
		return Snapshot{}, errors.New("label must not exceed 64 characters")
	}
	return Snapshot{
		ID:        uuid.NewString(),
		SnippetID: s.ID,
		Title:     s.Title,
		Content:   s.Content,
		Language:  s.Language,
		Tags:      append([]string(nil), s.Tags...),
		Label:     label,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// SnapshotsFor returns all snapshots belonging to the given snippet ID.
func SnapshotsFor(snippetID string, all []Snapshot) []Snapshot {
	var out []Snapshot
	for _, sn := range all {
		if sn.SnippetID == snippetID {
			out = append(out, sn)
		}
	}
	return out
}

// RemoveSnapshot removes the snapshot with the given ID from the slice.
func RemoveSnapshot(id string, all []Snapshot) ([]Snapshot, error) {
	for i, sn := range all {
		if sn.ID == id {
			return append(all[:i:i], all[i+1:]...), nil
		}
	}
	return nil, errors.New("snapshot not found")
}

// FindSnapshot returns the snapshot matching the given ID.
func FindSnapshot(id string, all []Snapshot) (Snapshot, bool) {
	for _, sn := range all {
		if sn.ID == id {
			return sn, true
		}
	}
	return Snapshot{}, false
}
