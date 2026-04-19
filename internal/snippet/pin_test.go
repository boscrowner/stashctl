package snippet

import (
	"testing"
	"time"
)

func makePinSnippet(id, title string, tags []string) Snippet {
	return Snippet{
		ID:        id,
		Title:     title,
		Content:   "content",
		Language:  "go",
		Tags:      tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestPinnedFilters(t *testing.T) {
	snippets := []Snippet{
		makePinSnippet("1", "Alpha", []string{"pinned", "go"}),
		makePinSnippet("2", "Beta", []string{"go"}),
		makePinSnippet("3", "Gamma", []string{"pinned"}),
	}
	result := Pinned(snippets)
	if len(result) != 2 {
		t.Fatalf("expected 2 pinned, got %d", len(result))
	}
	if result[0].Title != "Alpha" || result[1].Title != "Gamma" {
		t.Errorf("unexpected order: %v %v", result[0].Title, result[1].Title)
	}
}

func TestPinnedEmpty(t *testing.T) {
	result := Pinned([]Snippet{makePinSnippet("1", "A", []string{"go"})})
	if len(result) != 0 {
		t.Fatalf("expected 0, got %d", len(result))
	}
}

func TestPin(t *testing.T) {
	s := makePinSnippet("1", "A", []string{"go"})
	if !Pin(&s) {
		t.Fatal("expected Pin to return true")
	}
	if !hasTag(s, "pinned") {
		t.Error("expected pinned tag")
	}
	if Pin(&s) {
		t.Error("expected Pin to return false when already pinned")
	}
}

func TestUnpin(t *testing.T) {
	s := makePinSnippet("1", "A", []string{"go", "pinned"})
	if !Unpin(&s) {
		t.Fatal("expected Unpin to return true")
	}
	if hasTag(s, "pinned") {
		t.Error("expected pinned tag to be removed")
	}
	if Unpin(&s) {
		t.Error("expected Unpin to return false when not pinned")
	}
}
