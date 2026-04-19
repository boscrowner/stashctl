package main

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/stashctl/internal/snippet"
)

func TestRootIncludesRecentCmd(t *testing.T) {
	st, _ := setupTest(t)
	s := &snippet.Snippet{
		ID: "r1", Title: "RecentSnip", Content: "code", Language: "python",
		UpdatedAt: time.Now(),
	}
	_ = st.Add(s)

	cmd := newRecentCmd(st, false)
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("recent cmd failed: %v", err)
	}
	if !strings.Contains(buf.String(), "RecentSnip") {
		t.Errorf("expected RecentSnip in output: %s", buf.String())
	}
}
