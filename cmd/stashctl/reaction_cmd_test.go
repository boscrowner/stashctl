package main

import (
	"strings"
	"testing"
)

func TestReactionAddAndList(t *testing.T) {
	env := setupTest(t)

	// Add a snippet first
	run(t, env, "add", "--title", "My Snippet", "--content", "echo hi", "--language", "bash")

	// Retrieve its ID
	out := run(t, env, "list")
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		t.Fatal("expected at least one snippet")
	}
	snippetID := strings.Fields(lines[0])[0]

	// Add a reaction
	out = run(t, env, "reaction", "add", snippetID, "heart")
	if !strings.Contains(out, "heart") {
		t.Errorf("expected confirmation with emoji, got: %s", out)
	}

	// List reactions
	out = run(t, env, "reaction", "list", snippetID)
	if !strings.Contains(out, ":heart:") {
		t.Errorf("expected :heart: in output, got: %s", out)
	}
}

func TestReactionListEmpty(t *testing.T) {
	env := setupTest(t)
	out := run(t, env, "reaction", "list", "nonexistent-id")
	if !strings.Contains(out, "no reactions") {
		t.Errorf("expected 'no reactions', got: %s", out)
	}
}

func TestReactionRemoveNotFound(t *testing.T) {
	env := setupTest(t)
	err := runErr(t, env, "reaction", "remove", "bad-id")
	if err == nil {
		t.Error("expected error removing nonexistent reaction")
	}
}

func TestReactionUnknownEmoji(t *testing.T) {
	env := setupTest(t)
	run(t, env, "add", "--title", "S", "--content", "x", "--language", "bash")
	out := run(t, env, "list")
	snippetID := strings.Fields(strings.Split(strings.TrimSpace(out), "\n")[0])[0]

	err := runErr(t, env, "reaction", "add", snippetID, "unicorn")
	if err == nil {
		t.Error("expected error for unknown emoji")
	}
}
