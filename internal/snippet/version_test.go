package snippet_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/stashctl/internal/snippet"
)

func makeVersion(snippetID, content, message string) snippet.Version {
	v, err := snippet.NewVersion(snippetID, content, message)
	if err != nil {
		panic(err)
	}
	return v
}

func TestNewVersionValid(t *testing.T) {
	v, err := snippet.NewVersion("snip-1", "fmt.Println()", "initial")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.ID == "" {
		t.Error("expected non-empty ID")
	}
	if v.SnippetID != "snip-1" {
		t.Errorf("expected snip-1, got %s", v.SnippetID)
	}
	if v.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestNewVersionEmptySnippetID(t *testing.T) {
	_, err := snippet.NewVersion("", "content", "msg")
	if err == nil {
		t.Error("expected error for empty snippet id")
	}
}

func TestNewVersionEmptyContent(t *testing.T) {
	_, err := snippet.NewVersion("snip-1", "", "msg")
	if err == nil {
		t.Error("expected error for empty content")
	}
}

func TestNewVersionMessageTooLong(t *testing.T) {
	_, err := snippet.NewVersion("snip-1", "code", strings.Repeat("x", 201))
	if err == nil {
		t.Error("expected error for message too long")
	}
}

func TestVersionsFor(t *testing.T) {
	v1 := makeVersion("snip-1", "v1", "")
	v2 := makeVersion("snip-2", "v2", "")
	v3 := makeVersion("snip-1", "v3", "")
	all := []snippet.Version{v1, v2, v3}
	result := snippet.VersionsFor("snip-1", all)
	if len(result) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(result))
	}
}

func TestLatestVersion(t *testing.T) {
	v1 := makeVersion("snip-1", "old", "")
	time.Sleep(2 * time.Millisecond)
	v2 := makeVersion("snip-1", "new", "")
	all := []snippet.Version{v1, v2}
	latest, ok := snippet.LatestVersion("snip-1", all)
	if !ok {
		t.Fatal("expected a latest version")
	}
	if latest.Content != "new" {
		t.Errorf("expected 'new', got %s", latest.Content)
	}
}

func TestLatestVersionMissing(t *testing.T) {
	_, ok := snippet.LatestVersion("snip-99", []snippet.Version{})
	if ok {
		t.Error("expected no version found")
	}
}

func TestPruneVersions(t *testing.T) {
	var versions []snippet.Version
	for i := 0; i < 15; i++ {
		v := makeVersion("snip-1", "code", "")
		versions = append(versions, v)
	}
	pruned := snippet.PruneVersions("snip-1", versions, 10)
	if len(pruned) != 10 {
		t.Errorf("expected 10 versions after prune, got %d", len(pruned))
	}
}
