// Package snippet provides core types and operations for managing code snippets.
package snippet

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ChangeKind describes the type of change recorded in a changelog entry.
type ChangeKind string

const (
	// ChangeKindCreate represents the initial creation of a snippet.
	ChangeKindCreate ChangeKind = "create"
	// ChangeKindUpdate represents a content or metadata update.
	ChangeKindUpdate ChangeKind = "update"
	// ChangeKindTag represents a tag modification (add or remove).
	ChangeKindTag ChangeKind = "tag"
	// ChangeKindArchive represents an archive or unarchive action.
	ChangeKindArchive ChangeKind = "archive"
	// ChangeKindPin represents a pin or unpin action.
	ChangeKindPin ChangeKind = "pin"
	// ChangeKindFavorite represents a favorite or unfavorite action.
	ChangeKindFavorite ChangeKind = "favorite"
)

// ChangelogEntry records a single auditable change made to a snippet.
type ChangelogEntry struct {
	ID        string     `json:"id"`
	SnippetID string     `json:"snippet_id"`
	Kind      ChangeKind `json:"kind"`
	Summary   string     `json:"summary"`
	At        time.Time  `json:"at"`
}

// NewChangelogEntry creates a validated ChangelogEntry for the given snippet and change.
// Returns an error if snippetID or summary is empty, or if kind is unrecognised.
func NewChangelogEntry(snippetID string, kind ChangeKind, summary string) (ChangelogEntry, error) {
	if snippetID == "" {
		return ChangelogEntry{}, errors.New("changelog: snippetID must not be empty")
	}
	if summary == "" {
		return ChangelogEntry{}, errors.New("changelog: summary must not be empty")
	}
	if !isKnownChangeKind(kind) {
		return ChangelogEntry{}, errors.New("changelog: unknown change kind")
	}
	return ChangelogEntry{
		ID:        uuid.NewString(),
		SnippetID: snippetID,
		Kind:      kind,
		Summary:   summary,
		At:        time.Now().UTC(),
	}, nil
}

// ChangelogFor returns all entries that belong to the given snippetID,
// ordered from oldest to newest.
func ChangelogFor(snippetID string, entries []ChangelogEntry) []ChangelogEntry {
	var out []ChangelogEntry
	for _, e := range entries {
		if e.SnippetID == snippetID {
			out = append(out, e)
		}
	}
	// entries are appended in insertion order; sort ascending by time.
	sortChangelogAsc(out)
	return out
}

// LatestChange returns the most recent ChangelogEntry for snippetID.
// The second return value is false when no entries exist for that snippet.
func LatestChange(snippetID string, entries []ChangelogEntry) (ChangelogEntry, bool) {
	matches := ChangelogFor(snippetID, entries)
	if len(matches) == 0 {
		return ChangelogEntry{}, false
	}
	return matches[len(matches)-1], true
}

// PruneChangelog removes entries for snippetID, keeping only the most recent n.
// If n <= 0 all entries for that snippet are removed.
func PruneChangelog(snippetID string, entries []ChangelogEntry, n int) []ChangelogEntry {
	var kept, others []ChangelogEntry
	for _, e := range entries {
		if e.SnippetID == snippetID {
			kept = append(kept, e)
		} else {
			others = append(others, e)
		}
	}
	if n > 0 && len(kept) > n {
		kept = kept[len(kept)-n:]
	} else if n <= 0 {
		kept = nil
	}
	return append(others, kept...)
}

// isKnownChangeKind reports whether kind is one of the declared constants.
func isKnownChangeKind(kind ChangeKind) bool {
	switch kind {
	case ChangeKindCreate, ChangeKindUpdate, ChangeKindTag,
		ChangeKindArchive, ChangeKindPin, ChangeKindFavorite:
		return true
	}
	return false
}

// sortChangelogAsc sorts a slice of ChangelogEntry in ascending time order in-place.
func sortChangelogAsc(entries []ChangelogEntry) {
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].At.Before(entries[j-1].At); j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}
}
