package snippet

import (
	"strings"
	"testing"
)

func makeBadge(snippetID, name, color string) Badge {
	b, err := NewBadge(snippetID, name, "⭐", color, "alice", "")
	if err != nil {
		panic(err)
	}
	return b
}

func TestNewBadgeValid(t *testing.T) {
	b, err := NewBadge("s1", "Top Snippet", "⭐", "gold", "alice", "well done")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.ID == "" {
		t.Error("expected non-empty ID")
	}
	if b.Color != "gold" {
		t.Errorf("expected gold, got %s", b.Color)
	}
}

func TestNewBadgeDefaultColor(t *testing.T) {
	b, err := NewBadge("s1", "Nice", "", "", "bob", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Color != "gray" {
		t.Errorf("expected default gray, got %s", b.Color)
	}
}

func TestNewBadgeEmptySnippetID(t *testing.T) {
	_, err := NewBadge("", "Cool", "", "blue", "alice", "")
	if err == nil {
		t.Error("expected error for empty snippet ID")
	}
}

func TestNewBadgeEmptyName(t *testing.T) {
	_, err := NewBadge("s1", "", "", "blue", "alice", "")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestNewBadgeNameTooLong(t *testing.T) {
	_, err := NewBadge("s1", strings.Repeat("x", 41), "", "blue", "alice", "")
	if err == nil {
		t.Error("expected error for name too long")
	}
}

func TestNewBadgeUnknownColor(t *testing.T) {
	_, err := NewBadge("s1", "Cool", "", "neon-pink", "alice", "")
	if err == nil {
		t.Error("expected error for unknown color")
	}
}

func TestNewBadgeNoteTooLong(t *testing.T) {
	_, err := NewBadge("s1", "Cool", "", "blue", "alice", strings.Repeat("n", 201))
	if err == nil {
		t.Error("expected error for note too long")
	}
}

func TestBadgesFor(t *testing.T) {
	all := []Badge{
		makeBadge("s1", "First", "gold"),
		makeBadge("s2", "Second", "blue"),
		makeBadge("s1", "Third", "green"),
	}
	result := BadgesFor("s1", all)
	if len(result) != 2 {
		t.Errorf("expected 2, got %d", len(result))
	}
}

func TestRemoveBadge(t *testing.T) {
	b := makeBadge("s1", "Star", "gold")
	all := []Badge{b}
	out, removed := RemoveBadge(b.ID, all)
	if !removed {
		t.Error("expected removed=true")
	}
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d", len(out))
	}
}

func TestRemoveBadgeNotFound(t *testing.T) {
	b := makeBadge("s1", "Star", "gold")
	all := []Badge{b}
	_, removed := RemoveBadge("nonexistent", all)
	if removed {
		t.Error("expected removed=false")
	}
}

func TestFindBadge(t *testing.T) {
	b := makeBadge("s1", "Star", "gold")
	all := []Badge{b}
	found, ok := FindBadge(b.ID, all)
	if !ok {
		t.Fatal("expected to find badge")
	}
	if found.Name != "Star" {
		t.Errorf("unexpected name: %s", found.Name)
	}
}
