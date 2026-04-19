package snippet

import (
	"testing"
	"time"
)

var (
	t1 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 = time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	t3 = time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
)

func TestSortByTitleAsc(t *testing.T) {
	snippets := []Snippet{
		makeSortSnippet("zebra", "go", t1, t1),
		makeSortSnippet("apple", "go", t2, t2),
		makeSortSnippet("mango", "go", t3, t3),
	}
	Sort(snippets, SortOptions{Field: SortByTitle, Order: SortAsc})
	if snippets[0].Title != "apple" || snippets[2].Title != "zebra" {
		t.Errorf("unexpected order: %v", snippets)
	}
}

func TestSortByTitleDesc(t *testing.T) {
	snippets := []Snippet{
		makeSortSnippet("apple", "go", t1, t1),
		makeSortSnippet("zebra", "go", t2, t2),
	}
	Sort(snippets, SortOptions{Field: SortByTitle, Order: SortDesc})
	if snippets[0].Title != "zebra" {
		t.Errorf("expected zebra first, got %s", snippets[0].Title)
	}
}

func TestSortByCreatedAsc(t *testing.T) {
	snippets := []Snippet{
		makeSortSnippet("b", "go", t3, t1),
		makeSortSnippet("a", "go", t1, t1),
	}
	Sort(snippets, SortOptions{Field: SortByCreated, Order: SortAsc})
	if snippets[0].Title != "a" {
		t.Errorf("expected a first, got %s", snippets[0].Title)
	}
}

func TestSortByUpdatedDesc(t *testing.T) {
	snippets := []Snippet{
		makeSortSnippet("old", "go", t1, t1),
		makeSortSnippet("new", "go", t1, t3),
	}
	Sort(snippets, SortOptions{Field: SortByUpdated, Order: SortDesc})
	if snippets[0].Title != "new" {
		t.Errorf("expected new first, got %s", snippets[0].Title)
	}
}

func TestSortByLanguageAsc(t *testing.T) {
	snippets := []Snippet{
		makeSortSnippet("x", "python", t1, t1),
		makeSortSnippet("y", "go", t1, t1),
	}
	Sort(snippets, SortOptions{Field: SortByLanguage, Order: SortAsc})
	if snippets[0].Language != "go" {
		t.Errorf("expected go first, got %s", snippets[0].Language)
	}
}
