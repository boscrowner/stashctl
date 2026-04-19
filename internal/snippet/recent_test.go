package snippet

import (
	"testing"
	"time"
)

func makeRecentSnippet(title string, updatedAt time.Time) *Snippet {
	s := &Snippet{
		ID:        title,
		Title:     title,
		UpdatedAt: updatedAt,
	}
	return s
}

func TestRecentNoFilter(t *testing.T) {
	now := time.Now()
	snippets := []*Snippet{
		makeRecentSnippet("old", now.Add(-48*time.Hour)),
		makeRecentSnippet("new", now.Add(-1*time.Hour)),
		makeRecentSnippet("mid", now.Add(-24*time.Hour)),
	}
	result := Recent(snippets, RecentOptions{})
	if len(result) != 3 {
		t.Fatalf("expected 3, got %d", len(result))
	}
	if result[0].Title != "new" {
		t.Errorf("expected newest first, got %s", result[0].Title)
	}
}

func TestRecentLimit(t *testing.T) {
	now := time.Now()
	snippets := []*Snippet{
		makeRecentSnippet("a", now.Add(-1*time.Hour)),
		makeRecentSnippet("b", now.Add(-2*time.Hour)),
		makeRecentSnippet("c", now.Add(-3*time.Hour)),
	}
	result := Recent(snippets, RecentOptions{Limit: 2})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestRecentSinceFilter(t *testing.T) {
	now := time.Now()
	snippets := []*Snippet{
		makeRecentSnippet("recent", now.Add(-1*time.Hour)),
		makeRecentSnippet("old", now.Add(-72*time.Hour)),
	}
	result := Recent(snippets, RecentOptions{Since: now.Add(-48 * time.Hour)})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Title != "recent" {
		t.Errorf("unexpected snippet: %s", result[0].Title)
	}
}

func TestRecentEmpty(t *testing.T) {
	result := Recent(nil, RecentOptions{Limit: 5})
	if len(result) != 0 {
		t.Fatalf("expected empty, got %d", len(result))
	}
}
