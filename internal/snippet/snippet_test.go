package snippet

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	s := New("Hello", "fmt.Println()", "go", []string{"go", "print"})
	if s.ID == "" {
		t.Error("expected non-empty ID")
	}
	if s.Title != "Hello" {
		t.Errorf("unexpected title: %s", s.Title)
	}
	if s.Language != "go" {
		t.Errorf("unexpected language: %s", s.Language)
	}
	if len(s.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(s.Tags))
	}
	if s.CreatedAt.IsZero() || s.UpdatedAt.IsZero() {
		t.Error("expected timestamps to be set")
	}
}

func TestHasTag(t *testing.T) {
	s := New("t", "c", "go", []string{"alpha", "beta"})
	if !s.HasTag("alpha") {
		t.Error("expected HasTag alpha")
	}
	if s.HasTag("gamma") {
		t.Error("did not expect HasTag gamma")
	}
}

func TestHasAllTags(t *testing.T) {
	s := New("t", "c", "go", []string{"alpha", "beta"})
	if !s.HasAllTags([]string{"alpha", "beta"}) {
		t.Error("expected HasAllTags")
	}
	if s.HasAllTags([]string{"alpha", "gamma"}) {
		t.Error("did not expect HasAllTags with missing tag")
	}
}

func TestUpdate(t *testing.T) {
	s := New("old", "old content", "go", []string{"a"})
	origUpdated := s.UpdatedAt
	time.Sleep(2 * time.Millisecond)
	s.Update("new", "new content", "python", []string{"b", "c"})
	if s.Title != "new" {
		t.Errorf("unexpected title: %s", s.Title)
	}
	if s.Language != "python" {
		t.Errorf("unexpected language: %s", s.Language)
	}
	if !s.UpdatedAt.After(origUpdated) {
		t.Error("expected UpdatedAt to be bumped")
	}
	if len(s.Tags) != 2 {
		t.Errorf("expected 2 tags after update, got %d", len(s.Tags))
	}
}
