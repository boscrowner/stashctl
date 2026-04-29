package main

import (
	"strings"
	"testing"
)

func TestWebhookAddAndList(t *testing.T) {
	env := setupTest(t)

	// Add a snippet first
	env.args("add", "--title", "Hook Snippet", "--content", "echo hi", "--language", "bash")
	if err := env.cmd.Execute(); err != nil {
		t.Fatalf("add failed: %v", err)
	}

	snippets, err := env.store.List()
	if err != nil || len(snippets) == 0 {
		t.Fatal("expected at least one snippet")
	}
	sid := snippets[0].ID

	// Register webhook
	env.args("webhook", "add", sid, "https://example.com/hook", "--events", "created,updated")
	if err := env.cmd.Execute(); err != nil {
		t.Fatalf("webhook add failed: %v", err)
	}
	if !strings.Contains(env.out.String(), "webhook registered") {
		t.Errorf("expected registered message, got: %s", env.out.String())
	}

	// List webhooks
	env.args("webhook", "list", sid)
	if err := env.cmd.Execute(); err != nil {
		t.Fatalf("webhook list failed: %v", err)
	}
	if !strings.Contains(env.out.String(), "https://example.com/hook") {
		t.Errorf("expected URL in output, got: %s", env.out.String())
	}
}

func TestWebhookListEmpty(t *testing.T) {
	env := setupTest(t)
	env.args("webhook", "list", "nonexistent-snippet")
	if err := env.cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(env.out.String(), "no webhooks") {
		t.Errorf("expected empty message, got: %s", env.out.String())
	}
}

func TestWebhookRemove(t *testing.T) {
	env := setupTest(t)

	env.args("add", "--title", "T", "--content", "x", "--language", "go")
	_ = env.cmd.Execute()
	snippets, _ := env.store.List()
	sid := snippets[0].ID

	env.args("webhook", "add", sid, "https://hook.io/cb", "--events", "deleted")
	_ = env.cmd.Execute()

	hooks, _ := env.store.ListWebhooks(sid)
	if len(hooks) == 0 {
		t.Fatal("expected a webhook to be present")
	}

	env.args("webhook", "remove", hooks[0].ID)
	if err := env.cmd.Execute(); err != nil {
		t.Fatalf("remove failed: %v", err)
	}
	if !strings.Contains(env.out.String(), "removed") {
		t.Errorf("expected removed message, got: %s", env.out.String())
	}

	remaining, _ := env.store.ListWebhooks(sid)
	if len(remaining) != 0 {
		t.Errorf("expected 0 webhooks, got %d", len(remaining))
	}
}

func TestWebhookRemoveNotFound(t *testing.T) {
	env := setupTest(t)
	env.args("webhook", "remove", "ghost-id")
	if err := env.cmd.Execute(); err == nil {
		t.Error("expected error for missing webhook")
	}
}
