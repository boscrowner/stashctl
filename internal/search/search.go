package search

import (
	"strings"

	"github.com/user/stashctl/internal/snippet"
)

// Result holds a snippet and its relevance score.
type Result struct {
	Snippet snippet.Snippet
	Score   int
}

// ByQuery searches snippets by matching the query string against
// the snippet title, description, and body (case-insensitive).
func ByQuery(snippets []snippet.Snippet, query string) []Result {
	if query == "" {
		return nil
	}
	q := strings.ToLower(query)
	var results []Result
	for _, s := range snippets {
		score := scoreSnippet(s, q)
		if score > 0 {
			results = append(results, Result{Snippet: s, Score: score})
		}
	}
	sortResults(results)
	return results
}

// scoreSnippet returns a relevance score for the snippet against the query.
func scoreSnippet(s snippet.Snippet, query string) int {
	score := 0
	if strings.Contains(strings.ToLower(s.Title), query) {
		score += 3
	}
	if strings.Contains(strings.ToLower(s.Description), query) {
		score += 2
	}
	if strings.Contains(strings.ToLower(s.Body), query) {
		score += 1
	}
	return score
}

// sortResults sorts results descending by score (simple insertion sort).
func sortResults(results []Result) {
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Score > results[j-1].Score; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
}
