package snippet_test

import (
	"strings"
	"testing"

	"github.com/user/stashctl/internal/snippet"
)

func makeLink(t *testing.T, snippetID, rawURL, title string) snippet.Link {
	t.Helper()
	l, err := snippet.NewLink(snippetID, rawURL, title)
	if err != nil {
		t.Fatalf("makeLink: %v", err)
	}
	return l
}

func TestNewLinkValid(t *testing.T) {
	l, err := snippet.NewLink("s1", "https://example.com", "Example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.ID == "" {
		t.Error("expected non-empty ID")
	}
	if l.SnippetID != "s1" {
		t.Errorf("got snippet_id %q, want s1", l.SnippetID)
	}
}

func TestNewLinkEmptySnippetID(t *testing.T) {
	_, err := snippet.NewLink("", "https://example.com", "")
	if err == nil {
		t.Error("expected error for empty snippet id")
	}
}

func TestNewLinkInvalidURL(t *testing.T) {
	_, err := snippet.NewLink("s1", "not-a-url", "")
	if err == nil {
		t.Error("expected error for invalid url")
	}
}

func TestNewLinkNonHTTPScheme(t *testing.T) {
	_, err := snippet.NewLink("s1", "ftp://example.com/file", "")
	if err == nil {
		t.Error("expected error for ftp scheme")
	}
}

func TestNewLinkTitleTooLong(t *testing.T) {
	long := strings.Repeat("a", 121)
	_, err := snippet.NewLink("s1", "https://example.com", long)
	if err == nil {
		t.Error("expected error for title too long")
	}
}

func TestLinksFor(t *testing.T) {
	l1 := makeLink(t, "s1", "https://go.dev", "Go")
	l2 := makeLink(t, "s2", "https://rust-lang.org", "Rust")
	l3 := makeLink(t, "s1", "https://pkg.go.dev", "Pkg")

	all := []snippet.Link{l1, l2, l3}
	got := snippet.LinksFor(all, "s1")
	if len(got) != 2 {
		t.Fatalf("expected 2 links for s1, got %d", len(got))
	}
}

func TestRemoveLink(t *testing.T) {
	l1 := makeLink(t, "s1", "https://go.dev", "Go")
	l2 := makeLink(t, "s1", "https://pkg.go.dev", "Pkg")
	all := []snippet.Link{l1, l2}

	updated, ok := snippet.RemoveLink(all, l1.ID)
	if !ok {
		t.Fatal("expected removal to succeed")
	}
	if len(updated) != 1 {
		t.Fatalf("expected 1 link after removal, got %d", len(updated))
	}
	if updated[0].ID != l2.ID {
		t.Error("wrong link remained after removal")
	}
}

func TestRemoveLinkNotFound(t *testing.T) {
	l := makeLink(t, "s1", "https://go.dev", "")
	all := []snippet.Link{l}
	_, ok := snippet.RemoveLink(all, "nonexistent")
	if ok {
		t.Error("expected false for missing link")
	}
}

func TestFindLink(t *testing.T) {
	l := makeLink(t, "s1", "https://go.dev", "Go")
	all := []snippet.Link{l}

	found, ok := snippet.FindLink(all, l.ID)
	if !ok {
		t.Fatal("expected to find link")
	}
	if found.URL != "https://go.dev" {
		t.Errorf("got url %q, want https://go.dev", found.URL)
	}
}
