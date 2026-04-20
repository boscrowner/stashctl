package snippet

import (
	"strings"
	"testing"
)

func TestNewCollectionValid(t *testing.T) {
	c, err := NewCollection("My Snippets", "A handy collection")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Name != "My Snippets" {
		t.Errorf("expected name 'My Snippets', got %q", c.Name)
	}
	if c.Description != "A handy collection" {
		t.Errorf("unexpected description: %q", c.Description)
	}
	if c.ID == "" {
		t.Error("expected non-empty ID")
	}
	if len(c.SnippetIDs) != 0 {
		t.Errorf("expected empty snippet list, got %v", c.SnippetIDs)
	}
}

func TestNewCollectionEmptyName(t *testing.T) {
	_, err := NewCollection("", "desc")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestNewCollectionNameTooLong(t *testing.T) {
	_, err := NewCollection(strings.Repeat("x", 101), "")
	if err == nil {
		t.Fatal("expected error for name exceeding 100 chars")
	}
}

func TestAddSnippet(t *testing.T) {
	c, _ := NewCollection("Test", "")
	c.AddSnippet("id-1")
	c.AddSnippet("id-2")
	c.AddSnippet("id-1") // duplicate — should be ignored

	if len(c.SnippetIDs) != 2 {
		t.Errorf("expected 2 snippet IDs, got %d", len(c.SnippetIDs))
	}
}

func TestRemoveSnippet(t *testing.T) {
	c, _ := NewCollection("Test", "")
	c.AddSnippet("id-1")
	c.AddSnippet("id-2")

	removed := c.RemoveSnippet("id-1")
	if !removed {
		t.Error("expected RemoveSnippet to return true")
	}
	if len(c.SnippetIDs) != 1 || c.SnippetIDs[0] != "id-2" {
		t.Errorf("unexpected snippet list after remove: %v", c.SnippetIDs)
	}
}

func TestRemoveSnippetNotFound(t *testing.T) {
	c, _ := NewCollection("Test", "")
	removed := c.RemoveSnippet("nonexistent")
	if removed {
		t.Error("expected RemoveSnippet to return false for missing ID")
	}
}

func TestContains(t *testing.T) {
	c, _ := NewCollection("Test", "")
	c.AddSnippet("id-42")

	if !c.Contains("id-42") {
		t.Error("expected Contains to return true for added snippet")
	}
	if c.Contains("id-99") {
		t.Error("expected Contains to return false for absent snippet")
	}
}
