package snippet

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Version represents a saved snapshot of a snippet's content at a point in time.
type Version struct {
	ID        string    `json:"id"`
	SnippetID string    `json:"snippet_id"`
	Content   string    `json:"content"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// NewVersion creates a new Version snapshot for the given snippet.
func NewVersion(snippetID, content, message string) (Version, error) {
	if snippetID == "" {
		return Version{}, errors.New("snippet id must not be empty")
	}
	if content == "" {
		return Version{}, errors.New("content must not be empty")
	}
	if len(message) > 200 {
		return Version{}, errors.New("message must not exceed 200 characters")
	}
	return Version{
		ID:        uuid.NewString(),
		SnippetID: snippetID,
		Content:   content,
		Message:   message,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// VersionsFor returns all versions associated with the given snippet ID,
// ordered from oldest to newest.
func VersionsFor(snippetID string, versions []Version) []Version {
	var result []Version
	for _, v := range versions {
		if v.SnippetID == snippetID {
			result = append(result, v)
		}
	}
	return result
}

// LatestVersion returns the most recent version for a snippet, if any.
func LatestVersion(snippetID string, versions []Version) (Version, bool) {
	matches := VersionsFor(snippetID, versions)
	if len(matches) == 0 {
		return Version{}, false
	}
	return matches[len(matches)-1], true
}

// PruneVersions removes older versions, keeping at most maxKeep per snippet.
func PruneVersions(snippetID string, versions []Version, maxKeep int) []Version {
	if maxKeep <= 0 {
		maxKeep = 10
	}
	matches := VersionsFor(snippetID, versions)
	if len(matches) <= maxKeep {
		return versions
	}
	remove := make(map[string]bool)
	for _, v := range matches[:len(matches)-maxKeep] {
		remove[v.ID] = true
	}
	var pruned []Version
	for _, v := range versions {
		if !remove[v.ID] {
			pruned = append(pruned, v)
		}
	}
	return pruned
}
