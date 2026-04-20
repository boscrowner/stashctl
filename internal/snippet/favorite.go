package snippet

const favoriteTag = "favorite"

// Favorites returns snippets tagged as favorite.
func Favorites(snippets []*Snippet) []*Snippet {
	var out []*Snippet
	for _, s := range snippets {
		if hasTag(s, favoriteTag) {
			out = append(out, s)
		}
	}
	return out
}

// Favorite adds the "favorite" tag to a snippet if not already present.
func Favorite(s *Snippet) {
	for _, t := range s.Tags {
		if t == favoriteTag {
			return
		}
	}
	s.Tags = NormalizeTags(append(s.Tags, favoriteTag))
}

// Unfavorite removes the "favorite" tag from a snippet.
func Unfavorite(s *Snippet) {
	filtered := s.Tags[:0]
	for _, t := range s.Tags {
		if t != favoriteTag {
			filtered = append(filtered, t)
		}
	}
	s.Tags = filtered
}
