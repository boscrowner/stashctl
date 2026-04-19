package snippet

import "strings"

// DuplicateResult holds a pair of snippets considered duplicates.
type DuplicateResult struct {
	A *Snippet
	B *Snippet
	Score float64 // 0.0–1.0, higher means more similar
}

// FindDuplicates returns pairs of snippets whose titles and content are
// similar above the given threshold (0.0–1.0).
func FindDuplicates(snippets []*Snippet, threshold float64) []DuplicateResult {
	var results []DuplicateResult
	for i := 0; i < len(snippets); i++ {
		for j := i + 1; j < len(snippets); j++ {
			score := similarity(snippets[i], snippets[j])
			if score >= threshold {
				results = append(results, DuplicateResult{
					A:     snippets[i],
					B:     snippets[j],
					Score: score,
				})
			}
		}
	}
	return results
}

// similarity computes a simple similarity score between two snippets
// based on normalised title and content overlap.
func similarity(a, b *Snippet) float64 {
	titleScore := jaccardWords(a.Title, b.Title)
	contentScore := jaccardWords(a.Content, b.Content)
	return 0.4*titleScore + 0.6*contentScore
}

func jaccardWords(a, b string) float64 {
	setA := wordSet(a)
	setB := wordSet(b)
	if len(setA) == 0 && len(setB) == 0 {
		return 1.0
	}
	intersection := 0
	for w := range setA {
		if setB[w] {
			intersection++
		}
	}
	union := len(setA) + len(setB) - intersection
	if union == 0 {
		return 0
	}
	return float64(intersection) / float64(union)
}

func wordSet(s string) map[string]bool {
	set := make(map[string]bool)
	for _, w := range strings.Fields(strings.ToLower(s)) {
		set[w] = true
	}
	return set
}
