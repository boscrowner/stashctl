package snippet

import (
	"testing"
)

func TestNew(t *testing.T) {
	s := New("hello world", "fmt.Println(\"hello\")", "go", []string{"go", "print"})

	if s.Title != "hello world" {
		t.Errorf("expected title 'hello world', got %q", s.Title)
	}
	if s.Language != "go" {
		t.Errorf("expected language 'go', got %q", s.Language)
	}
	if len(s.ID) != 16 {
		t.Errorf("expected ID length 16, got %d", len(s.ID))
	}
	if s.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestHasTag(t *testing.T) {
	s := New("test", "content", "bash", []string{"shell", "util"})

	if !s.HasTag("shell") {
		t.Error("expected HasTag('shell') to be true")
	}
	if s.HasTag("python") {
		t.Error("expected HasTag('python') to be false")
	}
}

func TestHasAllTags(t *testing.T) {
	s := New("test", "content", "bash", []string{"shell", "util", "loop"})

	if !s.HasAllTags([]string{"shell", "loop"}) {
		t.Error("expected HasAllTags to return true for subset of tags")
	}
	if s.HasAllTags([]string{"shell", "missing"}) {
		t.Error("expected HasAllTags to return false when a tag is missing")
	}
	if !s.HasAllTags([]string{}) {
		t.Error("expected HasAllTags to return true for empty tag list")
	}
}
