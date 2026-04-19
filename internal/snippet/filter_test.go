package snippet

import (
	"testing"
)

func makeFilterSnippet(title, lang string, tags []string) *Snippet {
	s := &Snippet{
		ID:       "test-id",
		Title:    title,
		Language: lang,
		Tags:     tags,
	}
	return s
}

func TestByLanguage(t *testing.T) {
	snippets := []*Snippet{
		makeFilterSnippet("A", "go", nil),
		makeFilterSnippet("B", "python", nil),
		makeFilterSnippet("C", "go", nil),
	}
	got := ByLanguage(snippets, "go")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestByLanguageNoMatch(t *testing.T) {
	snippets := []*Snippet{
		makeFilterSnippet("A", "go", nil),
	}
	got := ByLanguage(snippets, "rust")
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestByTags(t *testing.T) {
	snippets := []*Snippet{
		makeFilterSnippet("A", "go", []string{"cli", "util"}),
		makeFilterSnippet("B", "go", []string{"cli"}),
		makeFilterSnippet("C", "go", []string{"util"}),
	}
	got := ByTags(snippets, []string{"cli", "util"})
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
	if got[0].Title != "A" {
		t.Errorf("expected snippet A, got %s", got[0].Title)
	}
}

func TestFilterApplyCombined(t *testing.T) {
	snippets := []*Snippet{
		makeFilterSnippet("A", "go", []string{"cli"}),
		makeFilterSnippet("B", "python", []string{"cli"}),
		makeFilterSnippet("C", "go", []string{"util"}),
	}
	f := Filter{Language: "go", Tags: []string{"cli"}}
	got := f.Apply(snippets)
	if len(got) != 1 || got[0].Title != "A" {
		t.Errorf("unexpected result: %v", got)
	}
}
