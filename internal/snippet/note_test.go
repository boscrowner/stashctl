package snippet

import (
	"strings"
	"testing"
)

func makeNote(snippetID, body string) Note {
	n, _ := NewNote(snippetID, body)
	return n
}

func TestNewNoteValid(t *testing.T) {
	n, err := NewNote("s1", "some note body")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.ID == "" {
		t.Error("expected non-empty ID")
	}
	if n.SnippetID != "s1" {
		t.Errorf("expected snippet_id s1, got %s", n.SnippetID)
	}
}

func TestNewNoteEmptySnippetID(t *testing.T) {
	_, err := NewNote("", "body")
	if err == nil {
		t.Error("expected error for empty snippet ID")
	}
}

func TestNewNoteEmptyBody(t *testing.T) {
	_, err := NewNote("s1", "")
	if err == nil {
		t.Error("expected error for empty body")
	}
}

func TestNewNoteBodyTooLong(t *testing.T) {
	_, err := NewNote("s1", strings.Repeat("x", 2001))
	if err == nil {
		t.Error("expected error for body exceeding limit")
	}
}

func TestNotesFor(t *testing.T) {
	notes := []Note{
		makeNote("s1", "note one"),
		makeNote("s2", "other"),
		makeNote("s1", "note two"),
	}
	got := NotesFor(notes, "s1")
	if len(got) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(got))
	}
}

func TestRemoveNote(t *testing.T) {
	n := makeNote("s1", "hello")
	notes := []Note{n}
	updated, ok := RemoveNote(notes, n.ID)
	if !ok {
		t.Error("expected removal to succeed")
	}
	if len(updated) != 0 {
		t.Errorf("expected empty slice, got %d items", len(updated))
	}
}

func TestRemoveNoteNotFound(t *testing.T) {
	notes := []Note{makeNote("s1", "hello")}
	_, ok := RemoveNote(notes, "nonexistent")
	if ok {
		t.Error("expected removal to fail for unknown id")
	}
}

func TestFindNote(t *testing.T) {
	n := makeNote("s1", "find me")
	notes := []Note{n}
	found, ok := FindNote(notes, n.ID)
	if !ok {
		t.Fatal("expected to find note")
	}
	if found.Body != "find me" {
		t.Errorf("unexpected body: %s", found.Body)
	}
}
