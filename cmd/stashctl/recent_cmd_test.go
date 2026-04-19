package main

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/stashctl/internal/snippet"
)

func TestRecentCmd(t *testing.T) {
	st, dir := setupTest(t)
	_ = dir

	s1 := &snippet.Snippet{
		ID: "id1", Title: "Alpha", Content: "a", Language: "go",
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}
	s2 := &snippet.Snippet{
		ID: "id2", Title: "Beta", Content: "b", Language: "go",
		UpdatedAt: time.Now().Add(-100 * 24 * time.Hour),
	}
	_ = st.Add(s1)
	_ = st.Add(s2)

	cmd := newRecentCmd(st, false)
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--since", "7"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Alpha") {
		t.Errorf("expected Alpha in output, got: %s", out)
	}
	if strings.Contains(out, "Beta") {
		t.Errorf("did not expect Beta in output, got: %s", out)
	}
}

func TestRecentCmdLimit(t *testing.T) {
	st, _ := setupTest(t)
	for i := 0; i < 5; i++ {
		s := &snippet.Snippet{
			ID: string(rune('a' + i)), Title: string(rune('A' + i)),
			Content: "x", Language: "go",
			UpdatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
		}
		_ = st.Add(s)
	}
	cmd := newRecentCmd(st, false)
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"-n", "2"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Count(buf.String(), "\n")
	if lines > 3 {
		t.Errorf("expected at most 2 result lines, got %d newlines", lines)
	}
}

func TestRecentCmdInvalidSince(t *testing.T) {
	st, _ := setupTest(t)
	cmd := newRecentCmd(st, false)
	cmd.SetArgs([]string{"--since", "notanumber"})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error for invalid --since value")
	}
}
