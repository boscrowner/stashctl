package snippet

import (
	"sort"
	"time"
)

// SortOrder defines the ordering direction.
type SortOrder int

const (
	SortAsc  SortOrder = iota
	SortDesc
)

// SortField defines which field to sort by.
type SortField string

const (
	SortByCreated  SortField = "created"
	SortByUpdated  SortField = "updated"
	SortByTitle    SortField = "title"
	SortByLanguage SortField = "language"
)

// SortOptions configures sorting behaviour.
type SortOptions struct {
	Field SortField
	Order SortOrder
}

// Sort sorts a slice of Snippets according to the provided options.
// Unknown fields fall back to sorting by created time descending.
func Sort(snippets []Snippet, opts SortOptions) {
	sort.SliceStable(snippets, func(i, j int) bool {
		var less bool
		switch opts.Field {
		case SortByUpdated:
			less = snippets[i].UpdatedAt.Before(snippets[j].UpdatedAt)
		case SortByTitle:
			less = snippets[i].Title < snippets[j].Title
		case SortByLanguage:
			less = snippets[i].Language < snippets[j].Language
		default: // SortByCreated
			less = snippets[i].CreatedAt.Before(snippets[j].CreatedAt)
		}
		if opts.Order == SortDesc {
			return !less
		}
		return less
	})
}

// makeSortSnippet is a helper used only within this package for testing.
func makeSortSnippet(title, lang string, created, updated time.Time) Snippet {
	return Snippet{
		Title:     title,
		Language:  lang,
		CreatedAt: created,
		UpdatedAt: updated,
	}
}
