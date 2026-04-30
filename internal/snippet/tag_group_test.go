package snippet

import (
	"testing"
)

func makeTagGroup(t *testing.T, name string, tags []string) TagGroup {
	t.Helper()
	g, err := NewTagGroup(name, "desc", tags)
	if err != nil {
		t.Fatalf("NewTagGroup: %v", err)
	}
	return g
}

func TestNewTagGroupValid(t *testing.T) {
	g, err := NewTagGroup("web", "web-related tags", []string{"http", "rest"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.ID == "" {
		t.Error("expected non-empty ID")
	}
	if g.Name != "web" {
		t.Errorf("expected name 'web', got %q", g.Name)
	}
	if len(g.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(g.Tags))
	}
}

func TestNewTagGroupEmptyName(t *testing.T) {
	_, err := NewTagGroup("", "", []string{"go"})
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestNewTagGroupNameTooLong(t *testing.T) {
	long := make([]byte, 65)
	for i := range long {
		long[i] = 'a'
	}
	_, err := NewTagGroup(string(long), "", []string{"go"})
	if err == nil {
		t.Error("expected error for name too long")
	}
}

func TestNewTagGroupNoTags(t *testing.T) {
	_, err := NewTagGroup("empty", "", []string{})
	if err == nil {
		t.Error("expected error for empty tags")
	}
}

func TestTagGroupsForSorted(t *testing.T) {
	g1 := makeTagGroup(t, "zebra", []string{"z"})
	g2 := makeTagGroup(t, "alpha", []string{"a"})
	result := TagGroupsFor([]TagGroup{g1, g2})
	if result[0].Name != "alpha" {
		t.Errorf("expected first group to be 'alpha', got %q", result[0].Name)
	}
}

func TestFindTagGroup(t *testing.T) {
	g := makeTagGroup(t, "backend", []string{"go", "api"})
	groups := []TagGroup{g}
	found, ok := FindTagGroup(groups, "Backend")
	if !ok {
		t.Fatal("expected to find tag group")
	}
	if found.ID != g.ID {
		t.Errorf("found wrong group: %v", found)
	}
}

func TestFindTagGroupMissing(t *testing.T) {
	g := makeTagGroup(t, "backend", []string{"go"})
	_, ok := FindTagGroup([]TagGroup{g}, "frontend")
	if ok {
		t.Error("expected not to find group")
	}
}

func TestRemoveTagGroup(t *testing.T) {
	g1 := makeTagGroup(t, "a", []string{"x"})
	g2 := makeTagGroup(t, "b", []string{"y"})
	result := RemoveTagGroup([]TagGroup{g1, g2}, g1.ID)
	if len(result) != 1 || result[0].ID != g2.ID {
		t.Errorf("unexpected result after remove: %v", result)
	}
}
