package snippet

import "sort"

// Pinned returns snippets marked as pinned (tag "pinned"), sorted by title.
func Pinned(snippets []Snippet) []Snippet {
	var result []Snippet
	for _, s := range snippets {
		if hasTag(s, "pinned") {
			result = append(result, s)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Title < result[j].Title
	})
	return result
}

// Pin adds the "pinned" tag to a snippet if not already present.
func Pin(s *Snippet) bool {
	if hasTag(*s, "pinned") {
		return false
	}
	s.Tags = NormalizeTags(append(s.Tags, "pinned"))
	return true
}

// Unpin removes the "pinned" tag from a snippet.
func Unpin(s *Snippet) bool {
	if !hasTag(*s, "pinned") {
		return false
	}
	filtered := s.Tags[:0]
	for _, t := range s.Tags {
		if t != "pinned" {
			filtered = append(filtered, t)
		}
	}
	s.Tags = filtered
	return true
}

func hasTag(s Snippet, tag string) bool {
	for _, t := range s.Tags {
		if t == tag {
			return true
		}
	}
	return false
}
