package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestWorkspaceCreateAndList(t *testing.T) {
	s, dir := setupTest(t)
	_ = dir

	root := newRootCmd(s)

	// create workspace
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"workspace", "create", "my-ws", "--description", "test desc"})
	if err := root.Execute(); err != nil {
		t.Fatalf("create workspace: %v", err)
	}
	if !strings.Contains(buf.String(), "my-ws") {
		t.Errorf("expected workspace name in output, got: %s", buf.String())
	}

	// list workspaces
	buf.Reset()
	root.SetArgs([]string{"workspace", "list"})
	if err := root.Execute(); err != nil {
		t.Fatalf("list workspaces: %v", err)
	}
	if !strings.Contains(buf.String(), "my-ws") {
		t.Errorf("expected workspace in list, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "test desc") {
		t.Errorf("expected description in list, got: %s", buf.String())
	}
}

func TestWorkspaceListEmpty(t *testing.T) {
	s, _ := setupTest(t)
	root := newRootCmd(s)
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"workspace", "list"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no workspaces") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWorkspaceAddAndRemoveSnippet(t *testing.T) {
	s, _ := setupTest(t)
	root := newRootCmd(s)

	// add a snippet first
	root.SetArgs([]string{"add", "--title", "hello", "--content", "echo hi", "--language", "bash"})
	if err := root.Execute(); err != nil {
		t.Fatalf("add snippet: %v", err)
	}
	snippets, _ := s.List()
	snippetID := snippets[0].ID

	// create workspace
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"workspace", "create", "ws1"})
	root.Execute()
	workspaces, _ := s.ListWorkspaces()
	wsID := workspaces[0].ID

	// add snippet to workspace
	buf.Reset()
	root.SetArgs([]string{"workspace", "add", wsID, snippetID})
	if err := root.Execute(); err != nil {
		t.Fatalf("workspace add: %v", err)
	}
	if !strings.Contains(buf.String(), "added") {
		t.Errorf("expected 'added' in output, got: %s", buf.String())
	}

	// remove snippet from workspace
	buf.Reset()
	root.SetArgs([]string{"workspace", "remove", wsID, snippetID})
	if err := root.Execute(); err != nil {
		t.Fatalf("workspace remove: %v", err)
	}
	if !strings.Contains(buf.String(), "removed") {
		t.Errorf("expected 'removed' in output, got: %s", buf.String())
	}
}

func TestWorkspaceCreateInvalidName(t *testing.T) {
	s, _ := setupTest(t)
	root := newRootCmd(s)
	root.SetArgs([]string{"workspace", "create", ""})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error for empty workspace name")
	}
}
