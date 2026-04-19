package snippet

import "time"

// Stats holds aggregate information about a collection of snippets.
type Stats struct {
	Total         int
	ByLanguage    map[string]int
	ByTag         map[string]int
	NewestUpdated time.Time
	OldestCreated time.Time
}

// ComputeStats returns statistics for the given snippet slice.
func ComputeStats(snippets []Snippet) Stats {
	s := Stats{
		ByLanguage: make(map[string]int),
		ByTag:      make(map[string]int),
	}
	for i, sn := range snippets {
		s.Total++
		if sn.Language != "" {
			s.ByLanguage[sn.Language]++
		}
		for _, t := range sn.Tags {
			s.ByTag[t]++
		}
		if i == 0 {
			s.OldestCreated = sn.CreatedAt
			s.NewestUpdated = sn.UpdatedAt
		} else {
			if sn.CreatedAt.Before(s.OldestCreated) {
				s.OldestCreated = sn.CreatedAt
			}
			if sn.UpdatedAt.After(s.NewestUpdated) {
				s.NewestUpdated = sn.UpdatedAt
			}
		}
	}
	return s
}
