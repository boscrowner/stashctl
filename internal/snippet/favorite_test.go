package snippet

import (
	"testing"
	"time"
)

func makeFavSnippet(id, title string, tags []string) *Snippet {
	return &Snippet{
		ID:        id,
		Title:     title,
		Content:   "content",
		Language:  "go",
		Tags:      tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestFavoritesFilters(t *testing.T) {
	snippets := []*Snippet{
		makeFavSnippet("1", "A", []string{"favorite", "go"}),
		makeFavSnippet("2", "B", []string{"go"}),
		makeFavSnippet("3", "C", []string{"favorite"}),
	}
	got := Favorites(snippets)
	if len(got) != 2 {
		t.Fatalf("expected 2 favorites, got %d", len(got))
	}
}

func TestFavoritesEmpty(t *testing.T) {
	got := Favorites([]*Snippet{})
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFavorite(t *testing.T) {
	s := makeFavSnippet("1", "A", []string{"go"})
	Favorite(s)
	if !hasTag(s, "favorite") {
		t.Fatal("expected favorite tag")
	}
	// idempotent
	Favorite(s)
	count := 0
	for _, tag := range s.Tags {
		if tag == "favorite" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("expected exactly 1 favorite tag, got %d", count)
	}
}

func TestUnfavorite(t *testing.T) {
	s := makeFavSnippet("1", "A", []string{"favorite", "go"})
	Unfavorite(s)
	if hasTag(s, "favorite") {
		t.Fatal("expected favorite tag removed")
	}
	if !hasTag(s, "go") {
		t.Fatal("expected go tag preserved")
	}
}
