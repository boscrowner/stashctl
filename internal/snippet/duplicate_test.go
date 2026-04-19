package snippet

import (
	"testing"
	"time"
)

func makeDupSnippet(id, title, content string) *Snippet {
	now := time.Now()
	return &Snippet{
		ID:        id,
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestFindDuplicatesNone(t *testing.T) {
	snippets := []*Snippet{
		makeDupSnippet("1", "hello world", "print hello"),
		makeDupSnippet("2", "fibonacci sequence", "func fib recursive"),
	}
	results := FindDuplicates(snippets, 0.5)
	if len(results) != 0 {
		t.Fatalf("expected no duplicates, got %d", len(results))
	}
}

func TestFindDuplicatesExact(t *testing.T) {
	s := makeDupSnippet("1", "hello world", "print hello world")
	clone := makeDupSnippet("2", "hello world", "print hello world")
	results := FindDuplicates([]*Snippet{s, clone}, 0.9)
	if len(results) != 1 {
		t.Fatalf("expected 1 duplicate pair, got %d", len(results))
	}
	if results[0].Score < 0.9 {
		t.Errorf("expected high score, got %f", results[0].Score)
	}
}

func TestFindDuplicatesThreshold(t *testing.T) {
	snippets := []*Snippet{
		makeDupSnippet("1", "sort a list", "sort items in a list using quicksort"),
		makeDupSnippet("2", "sort list items", "sort items in a list using mergesort"),
		makeDupSnippet("3", "connect to database", "open db connection pool"),
	}
	// low threshold should catch the two sort snippets
	results := FindDuplicates(snippets, 0.3)
	found := false
	for _, r := range results {
		if (r.A.ID == "1" && r.B.ID == "2") || (r.A.ID == "2" && r.B.ID == "1") {
			found = true
		}
	}
	if !found {
		t.Error("expected sort snippets to be flagged as duplicates")
	}
}

func TestJaccardWordsEmpty(t *testing.T) {
	if s := jaccardWords("", ""); s != 1.0 {
		t.Errorf("expected 1.0 for two empty strings, got %f", s)
	}
}
