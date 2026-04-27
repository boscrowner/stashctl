package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionSaveAndList(t *testing.T) {
	s, out := setupTest(t)

	// Add a snippet first
	out.Reset()
	root := newRootCmd(s)
	root.SetOut(out)
	root.SetArgs([]string{"add", "--title", "Hello", "--content", "fmt.Println()", "--language", "go"})
	if err := root.Execute(); err != nil {
		t.Fatalf("add failed: %v", err)
	}

	snippets, err := s.List()
	if err != nil || len(snippets) == 0 {
		t.Fatal("expected at least one snippet")
	}
	snipID := snippets[0].ID

	// Save a version
	out.Reset()
	root = newRootCmd(s)
	root.SetOut(out)
	root.SetArgs([]string{"version", "save", snipID, "--message", "first save"})
	if err := root.Execute(); err != nil {
		t.Fatalf("version save failed: %v", err)
	}
	if !strings.Contains(out.String(), "saved") {
		t.Errorf("expected 'saved' in output, got: %s", out.String())
	}

	// List versions
	out.Reset()
	root = newRootCmd(s)
	root.SetOut(out)
	root.SetArgs([]string{"version", "list", snipID})
	if err := root.Execute(); err != nil {
		t.Fatalf("version list failed: %v", err)
	}
	if !strings.Contains(out.String(), "first save") {
		t.Errorf("expected 'first save' in output, got: %s", out.String())
	}
}

func TestVersionListEmpty(t *testing.T) {
	s, out := setupTest(t)
	root := newRootCmd(s)
	root.SetOut(out)
	root.SetArgs([]string{"version", "list", "nonexistent-id"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "no versions found") {
		t.Errorf("expected 'no versions found', got: %s", out.String())
	}
}

func TestVersionSaveNotFound(t *testing.T) {
	s, _ := setupTest(t)
	root := newRootCmd(s)
	root.SetOut(&bytes.Buffer{})
	root.SetArgs([]string{"version", "save", "ghost-id"})
	if err := root.Execute(); err == nil {
		t.Error("expected error for unknown snippet id")
	}
}
