package snippet_test

import (
	"strings"
	"testing"

	"github.com/user/stashctl/internal/snippet"
)

func makeComment(snippetID, author, body string) snippet.Comment {
	c, err := snippet.NewComment(snippetID, author, body)
	if err != nil {
		panic(err)
	}
	return c
}

func TestNewCommentValid(t *testing.T) {
	c, err := snippet.NewComment("s1", "alice", "looks good")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.ID == "" {
		t.Error("expected non-empty ID")
	}
	if c.SnippetID != "s1" {
		t.Errorf("expected snippet_id s1, got %s", c.SnippetID)
	}
}

func TestNewCommentEmptySnippetID(t *testing.T) {
	_, err := snippet.NewComment("", "alice", "body")
	if err == nil {
		t.Error("expected error for empty snippet ID")
	}
}

func TestNewCommentEmptyAuthor(t *testing.T) {
	_, err := snippet.NewComment("s1", "", "body")
	if err == nil {
		t.Error("expected error for empty author")
	}
}

func TestNewCommentEmptyBody(t *testing.T) {
	_, err := snippet.NewComment("s1", "alice", "")
	if err == nil {
		t.Error("expected error for empty body")
	}
}

func TestNewCommentBodyTooLong(t *testing.T) {
	_, err := snippet.NewComment("s1", "alice", strings.Repeat("x", 1001))
	if err == nil {
		t.Error("expected error for body exceeding max length")
	}
}

func TestCommentsFor(t *testing.T) {
	c1 := makeComment("s1", "alice", "first")
	c2 := makeComment("s2", "bob", "second")
	c3 := makeComment("s1", "carol", "third")

	result := snippet.CommentsFor("s1", []snippet.Comment{c1, c2, c3})
	if len(result) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(result))
	}
}

func TestRemoveComment(t *testing.T) {
	c1 := makeComment("s1", "alice", "first")
	c2 := makeComment("s1", "bob", "second")

	updated, found := snippet.RemoveComment(c1.ID, []snippet.Comment{c1, c2})
	if !found {
		t.Error("expected comment to be found")
	}
	if len(updated) != 1 || updated[0].ID != c2.ID {
		t.Error("expected only c2 to remain")
	}
}

func TestRemoveCommentNotFound(t *testing.T) {
	c1 := makeComment("s1", "alice", "first")
	_, found := snippet.RemoveComment("nonexistent", []snippet.Comment{c1})
	if found {
		t.Error("expected not found")
	}
}

func TestFindComment(t *testing.T) {
	c1 := makeComment("s1", "alice", "hello")
	found, ok := snippet.FindComment(c1.ID, []snippet.Comment{c1})
	if !ok {
		t.Fatal("expected to find comment")
	}
	if found.Author != "alice" {
		t.Errorf("expected author alice, got %s", found.Author)
	}
}
