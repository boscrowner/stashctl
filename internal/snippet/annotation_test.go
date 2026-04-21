package snippet

import (
	"strings"
	"testing"
)

func TestNewAnnotationValid(t *testing.T) {
	a, err := NewAnnotation("snip-1", "This is useful for HTTP handlers.")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.ID == "" {
		t.Error("expected non-empty ID")
	}
	if a.SnippetID != "snip-1" {
		t.Errorf("expected snippet ID snip-1, got %s", a.SnippetID)
	}
	if a.Note != "This is useful for HTTP handlers." {
		t.Errorf("unexpected note: %s", a.Note)
	}
	if a.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestNewAnnotationEmptySnippetID(t *testing.T) {
	_, err := NewAnnotation("", "some note")
	if err == nil {
		t.Error("expected error for empty snippet ID")
	}
}

func TestNewAnnotationEmptyNote(t *testing.T) {
	_, err := NewAnnotation("snip-1", "")
	if err == nil {
		t.Error("expected error for empty note")
	}
}

func TestNewAnnotationNoteTooLong(t *testing.T) {
	long := strings.Repeat("x", 501)
	_, err := NewAnnotation("snip-1", long)
	if err == nil {
		t.Error("expected error for note exceeding 500 chars")
	}
}

func TestAnnotationsFor(t *testing.T) {
	anns := []Annotation{
		{ID: "a1", SnippetID: "s1", Note: "note1"},
		{ID: "a2", SnippetID: "s2", Note: "note2"},
		{ID: "a3", SnippetID: "s1", Note: "note3"},
	}
	result := AnnotationsFor(anns, "s1")
	if len(result) != 2 {
		t.Fatalf("expected 2 annotations, got %d", len(result))
	}
	for _, a := range result {
		if a.SnippetID != "s1" {
			t.Errorf("unexpected snippet ID: %s", a.SnippetID)
		}
	}
}

func TestRemoveAnnotation(t *testing.T) {
	anns := []Annotation{
		{ID: "a1", SnippetID: "s1", Note: "note1"},
		{ID: "a2", SnippetID: "s1", Note: "note2"},
	}
	updated, ok := RemoveAnnotation(anns, "a1")
	if !ok {
		t.Fatal("expected annotation to be removed")
	}
	if len(updated) != 1 {
		t.Fatalf("expected 1 annotation remaining, got %d", len(updated))
	}
	if updated[0].ID != "a2" {
		t.Errorf("expected remaining annotation ID a2, got %s", updated[0].ID)
	}
}

func TestRemoveAnnotationNotFound(t *testing.T) {
	anns := []Annotation{
		{ID: "a1", SnippetID: "s1", Note: "note1"},
	}
	updated, ok := RemoveAnnotation(anns, "nonexistent")
	if ok {
		t.Error("expected ok=false for missing annotation")
	}
	if len(updated) != 1 {
		t.Errorf("expected slice unchanged, got len %d", len(updated))
	}
}
