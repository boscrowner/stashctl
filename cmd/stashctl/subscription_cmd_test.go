package main

import (
	"strings"
	"testing"
)

func TestSubscriptionAddAndList(t *testing.T) {
	env := setupTest(t)

	// Add a snippet first
	run(t, env, "add", "--title", "My Snippet", "--content", "echo hi", "--lang", "bash")

	// Retrieve snippet ID from list
	out := run(t, env, "list")
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		t.Fatal("expected at least one snippet")
	}
	snippetID := strings.Fields(lines[0])[0]

	// Add subscription
	out = run(t, env, "subscription", "add", snippetID, "alice", "updated")
	if !strings.Contains(out, "subscription added") {
		t.Errorf("expected confirmation, got: %s", out)
	}

	// List subscriptions
	out = run(t, env, "subscription", "list", snippetID)
	if !strings.Contains(out, "alice") {
		t.Errorf("expected alice in subscription list, got: %s", out)
	}
	if !strings.Contains(out, "updated") {
		t.Errorf("expected event 'updated' in list, got: %s", out)
	}
}

func TestSubscriptionListEmpty(t *testing.T) {
	env := setupTest(t)
	out := run(t, env, "subscription", "list", "nonexistent")
	if !strings.Contains(out, "no subscriptions") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestSubscriptionRemove(t *testing.T) {
	env := setupTest(t)

	run(t, env, "add", "--title", "S", "--content", "x", "--lang", "go")
	out := run(t, env, "list")
	snippetID := strings.Fields(strings.Split(strings.TrimSpace(out), "\n")[0])[0]

	out = run(t, env, "subscription", "add", snippetID, "bob", "deleted")
	subID := strings.TrimPrefix(strings.TrimSpace(out), "subscription added: ")

	out = run(t, env, "subscription", "remove", subID)
	if !strings.Contains(out, "subscription removed") {
		t.Errorf("expected removed message, got: %s", out)
	}

	out = run(t, env, "subscription", "list", snippetID)
	if !strings.Contains(out, "no subscriptions") {
		t.Errorf("expected empty after removal, got: %s", out)
	}
}

func TestSubscriptionRemoveNotFound(t *testing.T) {
	env := setupTest(t)
	err := runErr(t, env, "subscription", "remove", "ghost-id")
	if err == nil {
		t.Error("expected error for unknown subscription ID")
	}
}
