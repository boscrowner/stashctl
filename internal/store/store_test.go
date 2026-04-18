package store_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/stashctl/internal/snippet"
	"github.com/user/stashctl/internal/store"
)

func tempStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	s, err := store.New(filepath.Join(dir, "snippets.json"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return s
}

func TestAddAndList(t *testing.T) {
	s := tempStore(t)
	sn := snippet.New("hello world", []string{"go", "example"})
	if err := s.Add(sn); err != nil {
		t.Fatalf("Add: %v", err)
	}
	list, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 snippet, got %d", len(list))
	}
	if list[0].Body != "hello world" {
		t.Errorf("unexpected body: %s", list[0].Body)
	}
}

func TestDelete(t *testing.T) {
	s := tempStore(t)
	sn := snippet.New("to delete", nil)
	_ = s.Add(sn)
	if err := s.Delete(sn.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	list, _ := s.List()
	if len(list) != 0 {
		t.Errorf("expected 0 snippets after delete, got %d", len(list))
	}
}

func TestDeleteNotFound(t *testing.T) {
	s := tempStore(t)
	if err := s.Delete("nonexistent"); err == nil {
		t.Error("expected error for missing ID, got nil")
	}
}

func TestFilterByTags(t *testing.T) {
	s := tempStore(t)
	_ = s.Add(snippet.New("a", []string{"go", "cli"}))
	_ = s.Add(snippet.New("b", []string{"go"}))
	_ = s.Add(snippet.New("c", []string{"python"}))

	results, err := s.FilterByTags([]string{"go"})
	if err != nil {
		t.Fatalf("FilterByTags: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	results, _ = s.FilterByTags([]string{"go", "cli"})
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestEmptyStore(t *testing.T) {
	s := tempStore(t)
	list, err := s.List()
	if err != nil {
		t.Fatalf("List on empty store: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d", len(list))
	}
	_ = os.Getenv("CI") // suppress unused import lint
}
