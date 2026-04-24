package snippet

import (
	"strings"
	"testing"
)

func makeLabel(t *testing.T, name, color string) Label {
	t.Helper()
	l, err := NewLabel(name, color)
	if err != nil {
		t.Fatalf("makeLabel: %v", err)
	}
	return l
}

func TestNewLabelValid(t *testing.T) {
	l, err := NewLabel("urgent", "#ff0000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.ID == "" {
		t.Error("expected non-empty ID")
	}
	if l.Name != "urgent" {
		t.Errorf("expected name 'urgent', got %q", l.Name)
	}
	if l.Color != "#ff0000" {
		t.Errorf("expected color '#ff0000', got %q", l.Color)
	}
}

func TestNewLabelDefaultColor(t *testing.T) {
	l, err := NewLabel("misc", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Color != "#cccccc" {
		t.Errorf("expected default color, got %q", l.Color)
	}
}

func TestNewLabelEmptyName(t *testing.T) {
	_, err := NewLabel("", "#fff")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestNewLabelNameTooLong(t *testing.T) {
	_, err := NewLabel(strings.Repeat("a", 33), "#fff")
	if err == nil {
		t.Error("expected error for name exceeding 32 chars")
	}
}

func TestLabelsFor(t *testing.T) {
	a := makeLabel(t, "alpha", "")
	b := makeLabel(t, "beta", "")
	c := makeLabel(t, "gamma", "")
	all := []Label{a, b, c}

	got := LabelsFor(all, []string{a.ID, c.ID})
	if len(got) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(got))
	}
	// sorted by name: alpha, gamma
	if got[0].Name != "alpha" || got[1].Name != "gamma" {
		t.Errorf("unexpected order: %v", got)
	}
}

func TestRemoveLabel(t *testing.T) {
	a := makeLabel(t, "alpha", "")
	b := makeLabel(t, "beta", "")
	all := []Label{a, b}

	result := RemoveLabel(all, a.ID)
	if len(result) != 1 || result[0].ID != b.ID {
		t.Errorf("expected only beta after removal, got %v", result)
	}
}

func TestFindLabel(t *testing.T) {
	a := makeLabel(t, "Alpha", "")
	all := []Label{a}

	l, ok := FindLabel(all, "alpha")
	if !ok {
		t.Fatal("expected to find label")
	}
	if l.ID != a.ID {
		t.Errorf("wrong label returned")
	}

	_, ok = FindLabel(all, "missing")
	if ok {
		t.Error("expected not found")
	}
}
