package snippet

import "time"

const tagArchived = "archived"

// Archived returns all snippets that have been archived.
func Archived(snippets []*Snippet) []*Snippet {
	var result []*Snippet
	for _, s := range snippets {
		if hasArchivedTag(s) {
			result = append(result, s)
		}
	}
	return result
}

// Active returns all snippets that have NOT been archived.
func Active(snippets []*Snippet) []*Snippet {
	var result []*Snippet
	for _, s := range snippets {
		if !hasArchivedTag(s) {
			result = append(result, s)
		}
	}
	return result
}

// Archive marks a snippet as archived by adding the "archived" tag
// and recording the current time as the updated timestamp.
func Archive(s *Snippet) {
	if hasArchivedTag(s) {
		return
	}
	s.Tags = NormalizeTags(append(s.Tags, tagArchived))
	s.UpdatedAt = time.Now()
}

// Unarchive removes the "archived" tag from a snippet.
func Unarchive(s *Snippet) {
	if !hasArchivedTag(s) {
		return
	}
	filtered := s.Tags[:0]
	for _, t := range s.Tags {
		if t != tagArchived {
			filtered = append(filtered, t)
		}
	}
	s.Tags = filtered
	s.UpdatedAt = time.Now()
}

func hasArchivedTag(s *Snippet) bool {
	for _, t := range s.Tags {
		if t == tagArchived {
			return true
		}
	}
	return false
}
