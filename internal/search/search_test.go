package search_test

import (
	"testing"

	"github.com/user/stashctl/internal/search"
	"github.com/user/stashctl/internal/snippet"
)

func makeSnippet(title, description, body string, tags []string) snippet.Snippet {
	return snippet.Snippet{
		Title:       title,
		Description: description,
		Body:        body,
		Tags:        tags,
	}
}

func TestByQueryEmpty(t *testing.T) {
	snippets := []snippet.Snippet{makeSnippet("hello", "world", "code", nil)}
	results := search.ByQuery(snippets, "")
	if len(results) != 0 {
		t.Fatalf("expected 0 results for empty query, got %d", len(results))
	}
}

func TestByQueryNoMatch(t *testing.T) {
	snippets := []snippet.Snippet{makeSnippet("hello", "world", "code", nil)}
	results := search.ByQuery(snippets, "xyz")
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestByQueryTitleMatch(t *testing.T) {
	snippets := []snippet.Snippet{
		makeSnippet("git stash", "save work", "git stash push", nil),
		makeSnippet("docker run", "run container", "docker run -it", nil),
	}
	results := search.ByQuery(snippets, "git")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Snippet.Title != "git stash" {
		t.Errorf("unexpected title: %s", results[0].Snippet.Title)
	}
}

func TestByQueryScoreOrdering(t *testing.T) {
	snippets := []snippet.Snippet{
		makeSnippet("unrelated", "unrelated", "contains curl", nil),
		makeSnippet("curl example", "curl usage", "curl http://example.com", nil),
	}
	results := search.ByQuery(snippets, "curl")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Snippet.Title != "curl example" {
		t.Errorf("expected highest scored result first, got: %s", results[0].Snippet.Title)
	}
}
