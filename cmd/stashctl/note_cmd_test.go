package main

import (
	"strings"
	"testing"
)

func TestNoteAddAndList(t *testing.T) {
	env := setupTest(t)

	// add a snippet first
	env.Execute(t, "add", "--title", "My Snippet", "--content", "echo hello", "--language", "bash")

	snippetID := env.FirstSnippetID(t)

	out := env.Execute(t, "note", "add", snippetID, "this is a note")
	if !strings.Contains(out, "added") {
		t.Errorf("expected 'added' in output, got: %s", out)
	}

	out = env.Execute(t, "note", "list", snippetID)
	if !strings.Contains(out, "this is a note") {
		t.Errorf("expected note body in output, got: %s", out)
	}
}

func TestNoteListEmpty(t *testing.T) {
	env := setupTest(t)
	env.Execute(t, "add", "--title", "S", "--content", "x", "--language", "text")
	snippetID := env.FirstSnippetID(t)

	out := env.Execute(t, "note", "list", snippetID)
	if !strings.Contains(out, "no notes") {
		t.Errorf("expected 'no notes', got: %s", out)
	}
}

func TestNoteRemove(t *testing.T) {
	env := setupTest(t)
	env.Execute(t, "add", "--title", "S", "--content", "x", "--language", "text")
	snippetID := env.FirstSnippetID(t)

	addOut := env.Execute(t, "note", "add", snippetID, "removable note")
	// extract note id prefix from "note <id> added"
	parts := strings.Fields(addOut)
	if len(parts) < 2 {
		t.Fatalf("unexpected add output: %s", addOut)
	}
	noteID := parts[1]

	out := env.Execute(t, "note", "remove", noteID)
	if !strings.Contains(out, "removed") {
		t.Errorf("expected 'removed' in output, got: %s", out)
	}

	listOut := env.Execute(t, "note", "list", snippetID)
	if strings.Contains(listOut, "removable note") {
		t.Error("expected note to be gone after removal")
	}
}

func TestNoteRemoveNotFound(t *testing.T) {
	env := setupTest(t)
	err := env.ExecuteErr(t, "note", "remove", "nonexistent-id")
	if err == nil {
		t.Error("expected error when removing non-existent note")
	}
}
