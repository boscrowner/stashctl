package snippet_test

import (
	"strings"
	"testing"

	"github.com/user/stashctl/internal/snippet"
)

func makeSnapshotSnippet(id, title, content string) snippet.Snippet {
	s, _ := snippet.New(title, content, "go", []string{"test"})
	s.ID = id
	return s
}

func TestNewSnapshotValid(t *testing.T) {
	s := makeSnapshotSnippet("s1", "Hello", "fmt.Println()")
	sn, err := snippet.NewSnapshot(s, "before-refactor")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sn.SnippetID != "s1" {
		t.Errorf("expected snippet_id s1, got %s", sn.SnippetID)
	}
	if sn.Label != "before-refactor" {
		t.Errorf("expected label before-refactor, got %s", sn.Label)
	}
	if sn.ID == "" {
		t.Error("expected non-empty snapshot id")
	}
}

func TestNewSnapshotEmptySnippetID(t *testing.T) {
	s := makeSnapshotSnippet("", "Hello", "body")
	_, err := snippet.NewSnapshot(s, "label")
	if err == nil {
		t.Error("expected error for empty snippet id")
	}
}

func TestNewSnapshotLabelTooLong(t *testing.T) {
	s := makeSnapshotSnippet("s1", "Hello", "body")
	_, err := snippet.NewSnapshot(s, strings.Repeat("x", 65))
	if err == nil {
		t.Error("expected error for label exceeding 64 chars")
	}
}

func TestSnapshotsFor(t *testing.T) {
	s1 := makeSnapshotSnippet("s1", "A", "a")
	s2 := makeSnapshotSnippet("s2", "B", "b")
	sn1, _ := snippet.NewSnapshot(s1, "")
	sn2, _ := snippet.NewSnapshot(s2, "")
	sn3, _ := snippet.NewSnapshot(s1, "v2")
	all := []snippet.Snapshot{sn1, sn2, sn3}
	result := snippet.SnapshotsFor("s1", all)
	if len(result) != 2 {
		t.Fatalf("expected 2 snapshots for s1, got %d", len(result))
	}
}

func TestRemoveSnapshot(t *testing.T) {
	s := makeSnapshotSnippet("s1", "A", "a")
	sn, _ := snippet.NewSnapshot(s, "")
	all := []snippet.Snapshot{sn}
	updated, err := snippet.RemoveSnapshot(sn.ID, all)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updated) != 0 {
		t.Errorf("expected empty slice after removal")
	}
}

func TestRemoveSnapshotNotFound(t *testing.T) {
	_, err := snippet.RemoveSnapshot("nonexistent", []snippet.Snapshot{})
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestFindSnapshot(t *testing.T) {
	s := makeSnapshotSnippet("s1", "A", "a")
	sn, _ := snippet.NewSnapshot(s, "tag")
	all := []snippet.Snapshot{sn}
	found, ok := snippet.FindSnapshot(sn.ID, all)
	if !ok {
		t.Fatal("expected snapshot to be found")
	}
	if found.Label != "tag" {
		t.Errorf("expected label tag, got %s", found.Label)
	}
}
