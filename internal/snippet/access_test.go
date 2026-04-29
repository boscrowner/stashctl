package snippet_test

import (
	"testing"

	"github.com/user/stashctl/internal/snippet"
)

func makeGrant(t *testing.T, snippetID, principal string, level snippet.AccessLevel) snippet.AccessGrant {
	t.Helper()
	g, err := snippet.NewAccessGrant(snippetID, principal, level, "")
	if err != nil {
		t.Fatalf("NewAccessGrant: %v", err)
	}
	return g
}

func TestNewAccessGrantValid(t *testing.T) {
	g, err := snippet.NewAccessGrant("s1", "alice", snippet.AccessRead, "for review")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.ID == "" {
		t.Error("expected non-empty ID")
	}
	if g.Level != snippet.AccessRead {
		t.Errorf("expected read, got %s", g.Level)
	}
}

func TestNewAccessGrantEmptySnippetID(t *testing.T) {
	_, err := snippet.NewAccessGrant("", "alice", snippet.AccessRead, "")
	if err == nil {
		t.Fatal("expected error for empty snippet ID")
	}
}

func TestNewAccessGrantEmptyPrincipal(t *testing.T) {
	_, err := snippet.NewAccessGrant("s1", "", snippet.AccessWrite, "")
	if err == nil {
		t.Fatal("expected error for empty principal")
	}
}

func TestNewAccessGrantUnknownLevel(t *testing.T) {
	_, err := snippet.NewAccessGrant("s1", "bob", snippet.AccessLevel("superuser"), "")
	if err == nil {
		t.Fatal("expected error for unknown access level")
	}
}

func TestNewAccessGrantNoteTooLong(t *testing.T) {
	long := make([]byte, 201)
	for i := range long {
		long[i] = 'x'
	}
	_, err := snippet.NewAccessGrant("s1", "bob", snippet.AccessRead, string(long))
	if err == nil {
		t.Fatal("expected error for note too long")
	}
}

func TestAccessGrantsFor(t *testing.T) {
	g1 := makeGrant(t, "s1", "alice", snippet.AccessRead)
	g2 := makeGrant(t, "s2", "bob", snippet.AccessWrite)
	g3 := makeGrant(t, "s1", "carol", snippet.AccessAdmin)

	result := snippet.AccessGrantsFor([]snippet.AccessGrant{g1, g2, g3}, "s1")
	if len(result) != 2 {
		t.Fatalf("expected 2 grants, got %d", len(result))
	}
}

func TestRemoveAccessGrant(t *testing.T) {
	g1 := makeGrant(t, "s1", "alice", snippet.AccessRead)
	g2 := makeGrant(t, "s1", "bob", snippet.AccessWrite)

	updated, err := snippet.RemoveAccessGrant([]snippet.AccessGrant{g1, g2}, g1.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updated) != 1 || updated[0].ID != g2.ID {
		t.Error("expected only g2 to remain")
	}
}

func TestRemoveAccessGrantNotFound(t *testing.T) {
	g := makeGrant(t, "s1", "alice", snippet.AccessRead)
	_, err := snippet.RemoveAccessGrant([]snippet.AccessGrant{g}, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing grant")
	}
}

func TestHasAccess(t *testing.T) {
	g := makeGrant(t, "s1", "alice", snippet.AccessWrite)
	grants := []snippet.AccessGrant{g}

	if !snippet.HasAccess(grants, "s1", "alice", snippet.AccessRead) {
		t.Error("write should satisfy read requirement")
	}
	if !snippet.HasAccess(grants, "s1", "alice", snippet.AccessWrite) {
		t.Error("write should satisfy write requirement")
	}
	if snippet.HasAccess(grants, "s1", "alice", snippet.AccessAdmin) {
		t.Error("write should not satisfy admin requirement")
	}
	if snippet.HasAccess(grants, "s1", "bob", snippet.AccessRead) {
		t.Error("bob has no grants")
	}
}
