package snippet

import (
	"testing"
	"time"
)

func makeHistorySnippet(id, content string) Snippet {
	return Snippet{ID: id, Title: "t", Content: content, Language: "go"}
}

func TestNewHistoryEntryValid(t *testing.T) {
	s := makeHistorySnippet("abc", "fmt.Println()")
	e, err := NewHistoryEntry(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.SnippetID != "abc" {
		t.Errorf("expected snippet_id abc, got %s", e.SnippetID)
	}
	if e.Content != "fmt.Println()" {
		t.Errorf("unexpected content: %s", e.Content)
	}
	if e.SavedAt.IsZero() {
		t.Error("expected SavedAt to be set")
	}
}

func TestNewHistoryEntryEmptyID(t *testing.T) {
	s := makeHistorySnippet("", "content")
	_, err := NewHistoryEntry(s)
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
}

func TestNewHistoryEntryEmptyContent(t *testing.T) {
	s := makeHistorySnippet("id1", "")
	_, err := NewHistoryEntry(s)
	if err == nil {
		t.Fatal("expected error for empty content")
	}
}

func TestHistoryForFilters(t *testing.T) {
	now := time.Now().UTC()
	entries := []HistoryEntry{
		{SnippetID: "a", Content: "v1", SavedAt: now.Add(-2 * time.Hour)},
		{SnippetID: "b", Content: "v1", SavedAt: now.Add(-1 * time.Hour)},
		{SnippetID: "a", Content: "v2", SavedAt: now},
	}
	result := HistoryFor("a", entries)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Content != "v1" || result[1].Content != "v2" {
		t.Error("entries not in chronological order")
	}
}

func TestLatestHistoryFound(t *testing.T) {
	now := time.Now().UTC()
	entries := []HistoryEntry{
		{SnippetID: "x", Content: "old", SavedAt: now.Add(-time.Hour)},
		{SnippetID: "x", Content: "new", SavedAt: now},
	}
	e, ok := LatestHistory("x", entries)
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if e.Content != "new" {
		t.Errorf("expected 'new', got %s", e.Content)
	}
}

func TestLatestHistoryNotFound(t *testing.T) {
	_, ok := LatestHistory("missing", []HistoryEntry{})
	if ok {
		t.Fatal("expected not found")
	}
}

func TestPruneHistory(t *testing.T) {
	now := time.Now().UTC()
	entries := []HistoryEntry{
		{SnippetID: "a", Content: "v1", SavedAt: now.Add(-3 * time.Hour)},
		{SnippetID: "a", Content: "v2", SavedAt: now.Add(-2 * time.Hour)},
		{SnippetID: "a", Content: "v3", SavedAt: now.Add(-1 * time.Hour)},
		{SnippetID: "b", Content: "v1", SavedAt: now},
	}
	pruned := PruneHistory(entries, 2)
	countA := 0
	for _, e := range pruned {
		if e.SnippetID == "a" {
			countA++
		}
	}
	if countA != 2 {
		t.Errorf("expected 2 entries for 'a', got %d", countA)
	}
	if len(pruned) != 3 {
		t.Errorf("expected 3 total entries, got %d", len(pruned))
	}
}

func TestPruneHistoryZeroMax(t *testing.T) {
	entries := []HistoryEntry{
		{SnippetID: "a", Content: "v1", SavedAt: time.Now()},
	}
	result := PruneHistory(entries, 0)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
