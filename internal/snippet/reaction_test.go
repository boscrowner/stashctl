package snippet_test

import (
	"testing"

	"github.com/user/stashctl/internal/snippet"
)

func makeReaction(t *testing.T, snippetID, actor, emoji string) snippet.Reaction {
	t.Helper()
	r, err := snippet.NewReaction(snippetID, actor, emoji)
	if err != nil {
		t.Fatalf("makeReaction: %v", err)
	}
	return r
}

func TestNewReactionValid(t *testing.T) {
	r, err := snippet.NewReaction("s1", "alice", "heart")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.ID == "" {
		t.Error("expected non-empty ID")
	}
	if r.Emoji != "heart" {
		t.Errorf("expected emoji heart, got %s", r.Emoji)
	}
}

func TestNewReactionEmptySnippetID(t *testing.T) {
	_, err := snippet.NewReaction("", "alice", "fire")
	if err == nil {
		t.Error("expected error for empty snippet_id")
	}
}

func TestNewReactionEmptyActor(t *testing.T) {
	_, err := snippet.NewReaction("s1", "", "fire")
	if err == nil {
		t.Error("expected error for empty actor")
	}
}

func TestNewReactionUnknownEmoji(t *testing.T) {
	_, err := snippet.NewReaction("s1", "alice", "unicorn")
	if err == nil {
		t.Error("expected error for unknown emoji")
	}
}

func TestReactionsFor(t *testing.T) {
	r1 := makeReaction(t, "s1", "alice", "thumbsup")
	r2 := makeReaction(t, "s2", "bob", "heart")
	r3 := makeReaction(t, "s1", "carol", "rocket")

	got := snippet.ReactionsFor([]snippet.Reaction{r1, r2, r3}, "s1")
	if len(got) != 2 {
		t.Fatalf("expected 2 reactions, got %d", len(got))
	}
}

func TestCountByEmoji(t *testing.T) {
	r1 := makeReaction(t, "s1", "alice", "thumbsup")
	r2 := makeReaction(t, "s1", "bob", "thumbsup")
	r3 := makeReaction(t, "s1", "carol", "heart")

	counts := snippet.CountByEmoji([]snippet.Reaction{r1, r2, r3})
	if counts["thumbsup"] != 2 {
		t.Errorf("expected 2 thumbsup, got %d", counts["thumbsup"])
	}
	if counts["heart"] != 1 {
		t.Errorf("expected 1 heart, got %d", counts["heart"])
	}
}

func TestRemoveReaction(t *testing.T) {
	r1 := makeReaction(t, "s1", "alice", "fire")
	r2 := makeReaction(t, "s1", "bob", "tada")

	result, err := snippet.RemoveReaction([]snippet.Reaction{r1, r2}, r1.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].ID != r2.ID {
		t.Error("expected only r2 to remain")
	}
}

func TestRemoveReactionNotFound(t *testing.T) {
	r1 := makeReaction(t, "s1", "alice", "eyes")
	_, err := snippet.RemoveReaction([]snippet.Reaction{r1}, "nonexistent")
	if err == nil {
		t.Error("expected error for missing reaction")
	}
}
