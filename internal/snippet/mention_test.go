package snippet_test

import (
	"strings"
	"testing"

	"github.com/user/stashctl/internal/snippet"
)

func makeMention(t *testing.T, from, to, ctx string) snippet.Mention {
	t.Helper()
	m, err := snippet.NewMention(from, to, ctx)
	if err != nil {
		t.Fatalf("NewMention: unexpected error: %v", err)
	}
	return m
}

func TestNewMentionValid(t *testing.T) {
	m := makeMention(t, "aaa", "bbb", "related implementation")
	if m.ID == "" {
		t.Error("expected non-empty ID")
	}
	if m.FromID != "aaa" || m.ToID != "bbb" {
		t.Errorf("unexpected IDs: %s → %s", m.FromID, m.ToID)
	}
}

func TestNewMentionEmptyFrom(t *testing.T) {
	_, err := snippet.NewMention("", "bbb", "")
	if err == nil {
		t.Error("expected error for empty from_id")
	}
}

func TestNewMentionEmptyTo(t *testing.T) {
	_, err := snippet.NewMention("aaa", "", "")
	if err == nil {
		t.Error("expected error for empty to_id")
	}
}

func TestNewMentionSelfLoop(t *testing.T) {
	_, err := snippet.NewMention("aaa", "aaa", "")
	if err == nil {
		t.Error("expected error for self-mention")
	}
}

func TestNewMentionContextTooLong(t *testing.T) {
	long := strings.Repeat("x", 281)
	_, err := snippet.NewMention("aaa", "bbb", long)
	if err == nil {
		t.Error("expected error for context exceeding 280 chars")
	}
}

func TestMentionsFrom(t *testing.T) {
	m1 := makeMention(t, "src", "dst1", "")
	m2 := makeMention(t, "src", "dst2", "")
	m3 := makeMention(t, "other", "dst1", "")

	result := snippet.MentionsFrom([]snippet.Mention{m1, m2, m3}, "src")
	if len(result) != 2 {
		t.Errorf("expected 2 mentions from src, got %d", len(result))
	}
}

func TestMentionsTo(t *testing.T) {
	m1 := makeMention(t, "src1", "dst", "")
	m2 := makeMention(t, "src2", "dst", "")
	m3 := makeMention(t, "src1", "other", "")

	result := snippet.MentionsTo([]snippet.Mention{m1, m2, m3}, "dst")
	if len(result) != 2 {
		t.Errorf("expected 2 mentions to dst, got %d", len(result))
	}
}

func TestRemoveMention(t *testing.T) {
	m1 := makeMention(t, "a", "b", "")
	m2 := makeMention(t, "a", "c", "")

	updated, err := snippet.RemoveMention([]snippet.Mention{m1, m2}, m1.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updated) != 1 || updated[0].ID != m2.ID {
		t.Error("expected only m2 to remain")
	}
}

func TestRemoveMentionNotFound(t *testing.T) {
	m := makeMention(t, "a", "b", "")
	_, err := snippet.RemoveMention([]snippet.Mention{m}, "nonexistent-id")
	if err == nil {
		t.Error("expected error when mention not found")
	}
}
