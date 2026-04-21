package snippet

import (
	"strings"
	"testing"
)

func makeBookmark(snippetID, name, note string) Bookmark {
	b, _ := NewBookmark(snippetID, name, note)
	return b
}

func TestNewBookmarkValid(t *testing.T) {
	b, err := NewBookmark("s1", "my-ref", "useful snippet")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.SnippetID != "s1" {
		t.Errorf("expected snippet ID s1, got %s", b.SnippetID)
	}
	if b.Name != "my-ref" {
		t.Errorf("expected name my-ref, got %s", b.Name)
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestNewBookmarkEmptyID(t *testing.T) {
	_, err := NewBookmark("", "name", "")
	if err == nil {
		t.Error("expected error for empty snippet ID")
	}
}

func TestNewBookmarkEmptyName(t *testing.T) {
	_, err := NewBookmark("s1", "", "")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestNewBookmarkNoteTooLong(t *testing.T) {
	long := strings.Repeat("x", 201)
	_, err := NewBookmark("s1", "ref", long)
	if err == nil {
		t.Error("expected error for note exceeding 200 chars")
	}
}

func TestBookmarksFor(t *testing.T) {
	all := []Bookmark{
		makeBookmark("s1", "ref-a", ""),
		makeBookmark("s2", "ref-b", ""),
		makeBookmark("s1", "ref-c", "note"),
	}
	result := BookmarksFor("s1", all)
	if len(result) != 2 {
		t.Fatalf("expected 2 bookmarks for s1, got %d", len(result))
	}
}

func TestBookmarksForEmpty(t *testing.T) {
	result := BookmarksFor("s1", nil)
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %d", len(result))
	}
}

func TestRemoveBookmark(t *testing.T) {
	all := []Bookmark{
		makeBookmark("s1", "ref-a", ""),
		makeBookmark("s1", "ref-b", ""),
	}
	updated, removed := RemoveBookmark("ref-a", "s1", all)
	if !removed {
		t.Error("expected removed to be true")
	}
	if len(updated) != 1 {
		t.Fatalf("expected 1 remaining bookmark, got %d", len(updated))
	}
	if updated[0].Name != "ref-b" {
		t.Errorf("expected ref-b to remain, got %s", updated[0].Name)
	}
}

func TestFindBookmark(t *testing.T) {
	all := []Bookmark{
		makeBookmark("s1", "ref-a", ""),
		makeBookmark("s2", "ref-b", ""),
	}
	b, ok := FindBookmark("ref-b", all)
	if !ok {
		t.Fatal("expected to find bookmark ref-b")
	}
	if b.SnippetID != "s2" {
		t.Errorf("expected snippet ID s2, got %s", b.SnippetID)
	}
}

func TestFindBookmarkMissing(t *testing.T) {
	_, ok := FindBookmark("nonexistent", []Bookmark{})
	if ok {
		t.Error("expected not found for missing bookmark")
	}
}
