package snippet_test

import (
	"strings"
	"testing"

	"github.com/user/stashctl/internal/snippet"
)

func makeWorkspace(t *testing.T, name string) *snippet.Workspace {
	t.Helper()
	w, err := snippet.NewWorkspace(name, "test workspace")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return w
}

func TestNewWorkspaceValid(t *testing.T) {
	w, err := snippet.NewWorkspace("my-workspace", "desc")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if w.Name != "my-workspace" {
		t.Errorf("expected name %q, got %q", "my-workspace", w.Name)
	}
	if w.ID == "" {
		t.Error("expected non-empty ID")
	}
	if len(w.SnippetIDs) != 0 {
		t.Error("expected empty snippet IDs")
	}
}

func TestNewWorkspaceEmptyName(t *testing.T) {
	_, err := snippet.NewWorkspace("", "desc")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestNewWorkspaceNameTooLong(t *testing.T) {
	_, err := snippet.NewWorkspace(strings.Repeat("a", 65), "desc")
	if err == nil {
		t.Fatal("expected error for name too long")
	}
}

func TestAddSnippetToWorkspace(t *testing.T) {
	w := makeWorkspace(t, "ws")
	if err := w.AddSnippet("snippet-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(w.SnippetIDs) != 1 || w.SnippetIDs[0] != "snippet-1" {
		t.Error("expected snippet-1 to be added")
	}
}

func TestAddSnippetIdempotent(t *testing.T) {
	w := makeWorkspace(t, "ws")
	w.AddSnippet("snippet-1")
	w.AddSnippet("snippet-1")
	if len(w.SnippetIDs) != 1 {
		t.Errorf("expected 1 snippet, got %d", len(w.SnippetIDs))
	}
}

func TestRemoveSnippetFromWorkspace(t *testing.T) {
	w := makeWorkspace(t, "ws")
	w.AddSnippet("snippet-1")
	if err := w.RemoveSnippet("snippet-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(w.SnippetIDs) != 0 {
		t.Error("expected empty snippet IDs after removal")
	}
}

func TestRemoveSnippetNotFound(t *testing.T) {
	w := makeWorkspace(t, "ws")
	if err := w.RemoveSnippet("missing"); err == nil {
		t.Fatal("expected error for missing snippet")
	}
}

func TestWorkspacesFor(t *testing.T) {
	w1 := makeWorkspace(t, "ws1")
	w2 := makeWorkspace(t, "ws2")
	w1.AddSnippet("s1")
	w2.AddSnippet("s1")
	w2.AddSnippet("s2")

	result := snippet.WorkspacesFor("s1", []*snippet.Workspace{w1, w2})
	if len(result) != 2 {
		t.Errorf("expected 2 workspaces, got %d", len(result))
	}

	result = snippet.WorkspacesFor("s2", []*snippet.Workspace{w1, w2})
	if len(result) != 1 {
		t.Errorf("expected 1 workspace, got %d", len(result))
	}
}
