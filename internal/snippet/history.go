package snippet

import (
	"errors"
	"time"
)

// HistoryEntry records a past version of a snippet's content.
type HistoryEntry struct {
	SnippetID string    `json:"snippet_id"`
	Content   string    `json:"content"`
	SavedAt   time.Time `json:"saved_at"`
}

// NewHistoryEntry creates a HistoryEntry capturing the current content of s.
func NewHistoryEntry(s Snippet) (HistoryEntry, error) {
	if s.ID == "" {
		return HistoryEntry{}, errors.New("snippet ID must not be empty")
	}
	if s.Content == "" {
		return HistoryEntry{}, errors.New("snippet content must not be empty")
	}
	return HistoryEntry{
		SnippetID: s.ID,
		Content:   s.Content,
		SavedAt:   time.Now().UTC(),
	}, nil
}

// HistoryFor returns all entries that belong to the given snippet ID,
// ordered from oldest to newest.
func HistoryFor(snippetID string, entries []HistoryEntry) []HistoryEntry {
	var result []HistoryEntry
	for _, e := range entries {
		if e.SnippetID == snippetID {
			result = append(result, e)
		}
	}
	// entries are assumed to be appended in chronological order;
	// sort defensively by SavedAt.
	for i := 1; i < len(result); i++ {
		for j := i; j > 0 && result[j].SavedAt.Before(result[j-1].SavedAt); j-- {
			result[j], result[j-1] = result[j-1], result[j]
		}
	}
	return result
}

// LatestHistory returns the most recent HistoryEntry for a snippet ID.
// It returns false if no entry exists.
func LatestHistory(snippetID string, entries []HistoryEntry) (HistoryEntry, bool) {
	all := HistoryFor(snippetID, entries)
	if len(all) == 0 {
		return HistoryEntry{}, false
	}
	return all[len(all)-1], true
}

// PruneHistory keeps only the most recent maxEntries entries per snippet,
// discarding older ones.
func PruneHistory(entries []HistoryEntry, maxEntries int) []HistoryEntry {
	if maxEntries <= 0 {
		return nil
	}
	// group by snippet ID preserving order
	index := make(map[string][]HistoryEntry)
	order := []string{}
	seen := make(map[string]bool)
	for _, e := range entries {
		if !seen[e.SnippetID] {
			order = append(order, e.SnippetID)
			seen[e.SnippetID] = true
		}
		index[e.SnippetID] = append(index[e.SnippetID], e)
	}
	var result []HistoryEntry
	for _, id := range order {
		group := index[id]
		if len(group) > maxEntries {
			group = group[len(group)-maxEntries:]
		}
		result = append(result, group...)
	}
	return result
}
