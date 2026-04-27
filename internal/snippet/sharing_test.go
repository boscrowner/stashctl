package snippet

import (
	"testing"
	"time"
)

func makeShareLink(t *testing.T, snippetID string, vis ShareVisibility, expiresAt *time.Time) ShareLink {
	t.Helper()
	sl, err := NewShareLink(snippetID, vis, expiresAt, "")
	if err != nil {
		t.Fatalf("NewShareLink: %v", err)
	}
	return sl
}

func TestNewShareLinkValid(t *testing.T) {
	sl := makeShareLink(t, "snip-1", VisibilityPublic, nil)
	if sl.SnippetID != "snip-1" {
		t.Errorf("expected snippet_id snip-1, got %s", sl.SnippetID)
	}
	if sl.Visibility != VisibilityPublic {
		t.Errorf("expected public visibility")
	}
	if sl.ID == "" || sl.Token == "" {
		t.Error("expected non-empty ID and Token")
	}
}

func TestNewShareLinkEmptySnippetID(t *testing.T) {
	_, err := NewShareLink("", VisibilityPublic, nil, "")
	if err == nil {
		t.Fatal("expected error for empty snippet ID")
	}
}

func TestNewShareLinkUnknownVisibility(t *testing.T) {
	_, err := NewShareLink("snip-1", "secret", nil, "")
	if err == nil {
		t.Fatal("expected error for unknown visibility")
	}
}

func TestNewShareLinkPastExpiry(t *testing.T) {
	past := time.Now().Add(-time.Hour)
	_, err := NewShareLink("snip-1", VisibilityPublic, &past, "")
	if err == nil {
		t.Fatal("expected error for past expiry")
	}
}

func TestNewShareLinkNoteTooLong(t *testing.T) {
	note := make([]byte, 201)
	for i := range note {
		note[i] = 'x'
	}
	_, err := NewShareLink("snip-1", VisibilityPublic, nil, string(note))
	if err == nil {
		t.Fatal("expected error for note exceeding 200 chars")
	}
}

func TestIsExpiredFalseWhenNoExpiry(t *testing.T) {
	sl := makeShareLink(t, "snip-1", VisibilityPrivate, nil)
	if sl.IsExpired() {
		t.Error("expected not expired when no expiry set")
	}
}

func TestIsExpiredTrueWhenPast(t *testing.T) {
	sl := makeShareLink(t, "snip-1", VisibilityPublic, nil)
	past := time.Now().Add(-time.Minute)
	sl.ExpiresAt = &past
	if !sl.IsExpired() {
		t.Error("expected expired for past expiry")
	}
}

func TestShareLinksFor(t *testing.T) {
	a := makeShareLink(t, "snip-1", VisibilityPublic, nil)
	b := makeShareLink(t, "snip-2", VisibilityUnlisted, nil)
	c := makeShareLink(t, "snip-1", VisibilityPrivate, nil)
	all := []ShareLink{a, b, c}
	result := ShareLinksFor(all, "snip-1")
	if len(result) != 2 {
		t.Fatalf("expected 2 links for snip-1, got %d", len(result))
	}
}

func TestRemoveShareLink(t *testing.T) {
	a := makeShareLink(t, "snip-1", VisibilityPublic, nil)
	b := makeShareLink(t, "snip-1", VisibilityPrivate, nil)
	all := []ShareLink{a, b}
	remaining, err := RemoveShareLink(all, a.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(remaining) != 1 || remaining[0].ID != b.ID {
		t.Error("expected only b to remain")
	}
}

func TestRemoveShareLinkNotFound(t *testing.T) {
	sl := makeShareLink(t, "snip-1", VisibilityPublic, nil)
	_, err := RemoveShareLink([]ShareLink{sl}, "nonexistent")
	if err == nil {
		t.Fatal("expected error when link not found")
	}
}
