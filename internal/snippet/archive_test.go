package snippet

import (
	"testing"
	"time"
)

func makeArchiveSnippet(id, title string, tags []string) *Snippet {
	return &Snippet{
		ID:        id,
		Title:     title,
		Content:   "echo hello",
		Language:  "bash",
		Tags:      tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestArchivedFilters(t *testing.T) {
	snippets := []*Snippet{
		makeArchiveSnippet("1", "Active One", []string{"go"}),
		makeArchiveSnippet("2", "Archived One", []string{"archived", "go"}),
		makeArchiveSnippet("3", "Archived Two", []string{"archived"}),
	}
	got := Archived(snippets)
	if len(got) != 2 {
		t.Fatalf("expected 2 archived snippets, got %d", len(got))
	}
}

func TestActiveFilters(t *testing.T) {
	snippets := []*Snippet{
		makeArchiveSnippet("1", "Active One", []string{"go"}),
		makeArchiveSnippet("2", "Archived One", []string{"archived"}),
	}
	got := Active(snippets)
	if len(got) != 1 {
		t.Fatalf("expected 1 active snippet, got %d", len(got))
	}
	if got[0].ID != "1" {
		t.Errorf("expected snippet id '1', got %s", got[0].ID)
	}
}

func TestArchive(t *testing.T) {
	s := makeArchiveSnippet("1", "My Snippet", []string{"go"})
	before := s.UpdatedAt
	time.Sleep(time.Millisecond)
	Archive(s)
	if !hasArchivedTag(s) {
		t.Error("expected snippet to have archived tag")
	}
	if !s.UpdatedAt.After(before) {
		t.Error("expected UpdatedAt to be refreshed after archive")
	}
}

func TestArchiveIdempotent(t *testing.T) {
	s := makeArchiveSnippet("1", "My Snippet", []string{"archived"})
	Archive(s)
	count := 0
	for _, t := range s.Tags {
		if t == "archived" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected exactly 1 archived tag, got %d", count)
	}
}

func TestUnarchive(t *testing.T) {
	s := makeArchiveSnippet("1", "My Snippet", []string{"archived", "go"})
	before := s.UpdatedAt
	time.Sleep(time.Millisecond)
	Unarchive(s)
	if hasArchivedTag(s) {
		t.Error("expected archived tag to be removed")
	}
	if !s.UpdatedAt.After(before) {
		t.Error("expected UpdatedAt to be refreshed after unarchive")
	}
}

func TestUnarchiveNoOp(t *testing.T) {
	s := makeArchiveSnippet("1", "My Snippet", []string{"go"})
	updated := s.UpdatedAt
	Unarchive(s)
	if s.UpdatedAt != updated {
		t.Error("expected UpdatedAt to remain unchanged when not archived")
	}
}
