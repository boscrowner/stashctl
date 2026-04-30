package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCommentAddAndList(t *testing.T) {
	env := setupTest(t)

	// add a snippet first
	env.Execute(t, "add", "--title", "my snippet", "--language", "go", "--content", "fmt.Println()")

	snippets := env.Execute(t, "list")
	lines := strings.Split(strings.TrimSpace(snippets), "\n")
	if len(lines) == 0 {
		t.Fatal("expected at least one snippet")
	}
	snippetID := strings.Fields(lines[0])[0]

	out := env.Execute(t, "comment", "add", snippetID, "great snippet", "--author", "alice")
	if !strings.Contains(out, "added") {
		t.Errorf("expected 'added' in output, got: %s", out)
	}

	list := env.Execute(t, "comment", "list", snippetID)
	if !strings.Contains(list, "alice") {
		t.Errorf("expected author alice in list, got: %s", list)
	}
	if !strings.Contains(list, "great snippet") {
		t.Errorf("expected body in list, got: %s", list)
	}
}

func TestCommentListEmpty(t *testing.T) {
	env := setupTest(t)
	env.Execute(t, "add", "--title", "s", "--language", "go", "--content", "x")
	snippets := env.Execute(t, "list")
	snippetID := strings.Fields(strings.TrimSpace(snippets))[0]

	out := env.Execute(t, "comment", "list", snippetID)
	if !strings.Contains(out, "no comments") {
		t.Errorf("expected 'no comments', got: %s", out)
	}
}

func TestCommentRemove(t *testing.T) {
	env := setupTest(t)
	env.Execute(t, "add", "--title", "s", "--language", "go", "--content", "x")
	snippets := env.Execute(t, "list")
	snippetID := strings.Fields(strings.TrimSpace(snippets))[0]

	addOut := env.Execute(t, "comment", "add", snippetID, "to remove", "--author", "bob")
	// extract comment id from output: "comment <id> added"
	parts := strings.Fields(addOut)
	var commentID string
	for i, p := range parts {
		if p == "comment" && i+1 < len(parts) {
			commentID = parts[i+1]
			break
		}
	}
	if commentID == "" {
		t.Fatal("could not parse comment ID from output")
	}

	out := env.Execute(t, "comment", "remove", commentID)
	if !strings.Contains(out, "removed") {
		t.Errorf("expected 'removed' in output, got: %s", out)
	}

	list := env.Execute(t, "comment", "list", snippetID)
	if strings.Contains(list, "bob") {
		t.Error("expected comment to be removed")
	}
}

func TestCommentRemoveNotFound(t *testing.T) {
	env := setupTest(t)
	var buf bytes.Buffer
	env.Cmd.SetErr(&buf)
	env.Cmd.SetArgs([]string{"comment", "remove", "nonexistent-id"})
	err := env.Cmd.Execute()
	if err == nil {
		t.Error("expected error when removing nonexistent comment")
	}
}
