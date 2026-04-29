package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestReviewAddAndList(t *testing.T) {
	s, dir := setupTest(t)
	_ = dir

	// Add a snippet first
	root := newRootCmd(s)
	root.SetArgs([]string{"add", "--title", "My Snippet", "--content", "echo hi", "--language", "bash"})
	if err := root.Execute(); err != nil {
		t.Fatalf("add snippet: %v", err)
	}

	snippets, _ := s.List()
	if len(snippets) == 0 {
		t.Fatal("expected at least one snippet")
	}
	snipID := snippets[0].ID

	// Add a review
	var buf bytes.Buffer
	root = newRootCmd(s)
	root.SetOut(&buf)
	root.SetArgs([]string{"review", "add", snipID, "alice", "--status", "approved", "--comment", "LGTM"})
	if err := root.Execute(); err != nil {
		t.Fatalf("review add: %v", err)
	}
	if !strings.Contains(buf.String(), "added") {
		t.Errorf("expected 'added' in output, got: %s", buf.String())
	}

	// List reviews
	buf.Reset()
	root = newRootCmd(s)
	root.SetOut(&buf)
	root.SetArgs([]string{"review", "list", snipID})
	if err := root.Execute(); err != nil {
		t.Fatalf("review list: %v", err)
	}
	if !strings.Contains(buf.String(), "alice") {
		t.Errorf("expected reviewer 'alice' in output, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "approved") {
		t.Errorf("expected status 'approved' in output, got: %s", buf.String())
	}
}

func TestReviewListEmpty(t *testing.T) {
	s, _ := setupTest(t)
	var buf bytes.Buffer
	root := newRootCmd(s)
	root.SetOut(&buf)
	root.SetArgs([]string{"review", "list", "nonexistent"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no reviews found") {
		t.Errorf("expected 'no reviews found', got: %s", buf.String())
	}
}

func TestReviewRemoveNotFound(t *testing.T) {
	s, _ := setupTest(t)
	root := newRootCmd(s)
	root.SetArgs([]string{"review", "remove", "does-not-exist"})
	if err := root.Execute(); err == nil {
		t.Error("expected error removing non-existent review")
	}
}
