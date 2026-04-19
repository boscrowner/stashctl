package snippet

// Filter holds criteria for filtering snippets.
type Filter struct {
	Tags     []string
	Language string
}

// Snippet is referenced from snippet.go; Filter operates on slices of *Snippet.

// Apply returns only those snippets matching all criteria in f.
func (f Filter) Apply(snippets []*Snippet) []*Snippet {
	var out []*Snippet
	for _, s := range snippets {
		if f.Language != "" && s.Language != f.Language {
			continue
		}
		if len(f.Tags) > 0 && !s.HasAllTags(f.Tags) {
			continue
		}
		out = append(out, s)
	}
	return out
}

// ByLanguage returns snippets whose Language field matches lang (case-sensitive).
func ByLanguage(snippets []*Snippet, lang string) []*Snippet {
	return Filter{Language: lang}.Apply(snippets)
}

// ByTags returns snippets that have all of the given tags.
func ByTags(snippets []*Snippet, tags []string) []*Snippet {
	return Filter{Tags: tags}.Apply(snippets)
}
