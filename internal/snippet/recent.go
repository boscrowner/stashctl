package snippet

import (
	"sort"
	"time"
)

// RecentOptions configures the Recent filter.
type RecentOptions struct {
	Limit  int
	Since  time.Time
}

// Recent returns snippets updated within the given options.
// If Limit <= 0, all matching snippets are returned.
// If Since is zero, no time filter is applied.
// Results are sorted by UpdatedAt descending (most recent first).
func Recent(snippets []*Snippet, opts RecentOptions) []*Snippet {
	var filtered []*Snippet
	for _, s := range snippets {
		if !opts.Since.IsZero() && s.UpdatedAt.Before(opts.Since) {
			continue
		}
		filtered = append(filtered, s)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].UpdatedAt.After(filtered[j].UpdatedAt)
	})

	if opts.Limit > 0 && len(filtered) > opts.Limit {
		return filtered[:opts.Limit]
	}
	return filtered
}

// RecentN is a convenience wrapper around Recent that returns at most n
// snippets updated since the given duration ago (e.g. 24*time.Hour).
func RecentN(snippets []*Snippet, n int, age time.Duration) []*Snippet {
	opts := RecentOptions{
		Limit: n,
	}
	if age > 0 {
		opts.Since = time.Now().Add(-age)
	}
	return Recent(snippets, opts)
}
