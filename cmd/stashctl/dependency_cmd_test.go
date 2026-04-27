package main

import (
	"strings"
	"testing"
)

func TestDepAddAndList(t *testing.T) {
	env := setupTest(t)

	// Add two snippets first
	env.run(t, "add", "--title", "Alpha", "--content", "alpha code", "--lang", "go")
	env.run(t, "add", "--title", "Beta", "--content", "beta code", "--lang", "go")

	snippets := env.listSnippets(t)
	if len(snippets) < 2 {
		t.Fatalf("expected at least 2 snippets, got %d", len(snippets))
	}
	srcID := snippets[0].ID
	tgtID := snippets[1].ID

	// Add dependency
	out := env.run(t, "dep", "add", srcID, tgtID, "--note", "alpha needs beta")
	if !strings.Contains(out, "dependency added") {
		t.Errorf("expected confirmation, got: %s", out)
	}

	// List dependencies for source
	out = env.run(t, "dep", "list", srcID)
	if !strings.Contains(out, tgtID) {
		t.Errorf("expected target ID in output, got: %s", out)
	}
	if !strings.Contains(out, "alpha needs beta") {
		t.Errorf("expected note in output, got: %s", out)
	}
}

func TestDepListDependents(t *testing.T) {
	env := setupTest(t)

	env.run(t, "add", "--title", "A", "--content", "a", "--lang", "go")
	env.run(t, "add", "--title", "B", "--content", "b", "--lang", "go")
	snippets := env.listSnippets(t)
	srcID := snippets[0].ID
	tgtID := snippets[1].ID

	env.run(t, "dep", "add", srcID, tgtID)

	out := env.run(t, "dep", "list", "--dependents", tgtID)
	if !strings.Contains(out, srcID) {
		t.Errorf("expected source in dependents output, got: %s", out)
	}
}

func TestDepRemove(t *testing.T) {
	env := setupTest(t)

	env.run(t, "add", "--title", "X", "--content", "x", "--lang", "go")
	env.run(t, "add", "--title", "Y", "--content", "y", "--lang", "go")
	snippets := env.listSnippets(t)
	srcID := snippets[0].ID
	tgtID := snippets[1].ID

	out := env.run(t, "dep", "add", srcID, tgtID)
	// Extract dependency ID from output: "dependency added: ... (id: <id>)"
	parts := strings.Split(out, "id: ")
	if len(parts) < 2 {
		t.Fatalf("could not parse dep ID from: %s", out)
	}
	depID := strings.TrimRight(strings.TrimSpace(parts[1]), ")")

	out = env.run(t, "dep", "remove", depID)
	if !strings.Contains(out, "removed") {
		t.Errorf("expected removed confirmation, got: %s", out)
	}

	out = env.run(t, "dep", "list", srcID)
	if !strings.Contains(out, "no dependencies found") {
		t.Errorf("expected empty list after removal, got: %s", out)
	}
}

func TestDepRemoveNotFound(t *testing.T) {
	env := setupTest(t)
	_, err := env.runErr("dep", "remove", "nonexistent-id")
	if err == nil {
		t.Error("expected error removing nonexistent dependency")
	}
}
