package snippet

import (
	"strings"
	"testing"
)

func makeAlias(t *testing.T, name, snippetID, note string) Alias {
	t.Helper()
	a, err := NewAlias(name, snippetID, note)
	if err != nil {
		t.Fatalf("makeAlias: %v", err)
	}
	return a
}

func TestNewAliasValid(t *testing.T) {
	a, err := NewAlias("mysnip", "snip-001", "quick access")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Name != "mysnip" {
		t.Errorf("expected name mysnip, got %s", a.Name)
	}
	if a.SnippetID != "snip-001" {
		t.Errorf("expected snippet ID snip-001, got %s", a.SnippetID)
	}
	if a.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestNewAliasEmptyName(t *testing.T) {
	_, err := NewAlias("", "snip-001", "")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestNewAliasWhitespaceName(t *testing.T) {
	_, err := NewAlias("my alias", "snip-001", "")
	if err == nil {
		t.Error("expected error for name with whitespace")
	}
}

func TestNewAliasNameTooLong(t *testing.T) {
	long := strings.Repeat("a", maxAliasNameLen+1)
	_, err := NewAlias(long, "snip-001", "")
	if err == nil {
		t.Error("expected error for name exceeding max length")
	}
}

func TestNewAliasEmptySnippetID(t *testing.T) {
	_, err := NewAlias("mysnip", "", "")
	if err == nil {
		t.Error("expected error for empty snippet ID")
	}
}

func TestAliasesFor(t *testing.T) {
	a1 := makeAlias(t, "foo", "snip-001", "")
	a2 := makeAlias(t, "bar", "snip-002", "")
	a3 := makeAlias(t, "baz", "snip-001", "")

	result := AliasesFor([]Alias{a1, a2, a3}, "snip-001")
	if len(result) != 2 {
		t.Fatalf("expected 2 aliases, got %d", len(result))
	}
}

func TestFindAlias(t *testing.T) {
	a := makeAlias(t, "mysnip", "snip-001", "")
	aliases := []Alias{a}

	found, ok := FindAlias(aliases, "mysnip")
	if !ok {
		t.Fatal("expected alias to be found")
	}
	if found.SnippetID != "snip-001" {
		t.Errorf("expected snip-001, got %s", found.SnippetID)
	}

	_, ok = FindAlias(aliases, "missing")
	if ok {
		t.Error("expected alias not to be found")
	}
}

func TestRemoveAlias(t *testing.T) {
	a1 := makeAlias(t, "foo", "snip-001", "")
	a2 := makeAlias(t, "bar", "snip-002", "")

	result, err := RemoveAlias([]Alias{a1, a2}, a1.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 alias remaining, got %d", len(result))
	}
	if result[0].ID != a2.ID {
		t.Error("wrong alias removed")
	}
}

func TestRemoveAliasNotFound(t *testing.T) {
	a := makeAlias(t, "foo", "snip-001", "")
	_, err := RemoveAlias([]Alias{a}, "nonexistent-id")
	if err == nil {
		t.Error("expected error when alias not found")
	}
}
