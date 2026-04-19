package snippet_test

import (
	"testing"
	"time"

	"github.com/user/stashctl/internal/snippet"
)

func makeStatsSnippet(title, lang string, tags []string, created, updated time.Time) snippet.Snippet {
	return snippet.Snippet{
		ID:        title,
		Title:     title,
		Language:  lang,
		Tags:      tags,
		CreatedAt: created,
		UpdatedAt: updated,
	}
}

func TestComputeStatsEmpty(t *testing.T) {
	s := snippet.ComputeStats(nil)
	if s.Total != 0 {
		t.Fatalf("expected 0, got %d", s.Total)
	}
}

func TestComputeStatsTotal(t *testing.T) {
	now := time.Now()
	snippets := []snippet.Snippet{
		makeStatsSnippet("a", "go", []string{"cli"}, now, now),
		makeStatsSnippet("b", "go", []string{"cli", "util"}, now, now),
		makeStatsSnippet("c", "python", []string{"util"}, now, now),
	}
	s := snippet.ComputeStats(snippets)
	if s.Total != 3 {
		t.Fatalf("expected 3, got %d", s.Total)
	}
	if s.ByLanguage["go"] != 2 {
		t.Fatalf("expected 2 go snippets, got %d", s.ByLanguage["go"])
	}
	if s.ByLanguage["python"] != 1 {
		t.Fatalf("expected 1 python snippet")
	}
	if s.ByTag["cli"] != 2 {
		t.Fatalf("expected 2 cli tags, got %d", s.ByTag["cli"])
	}
	if s.ByTag["util"] != 2 {
		t.Fatalf("expected 2 util tags")
	}
}

func TestComputeStatsDates(t *testing.T) {
	old := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	new_ := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	snippets := []snippet.Snippet{
		makeStatsSnippet("a", "go", nil, old, new_),
		makeStatsSnippet("b", "go", nil, new_, old),
	}
	s := snippet.ComputeStats(snippets)
	if !s.OldestCreated.Equal(old) {
		t.Fatalf("expected oldest %v, got %v", old, s.OldestCreated)
	}
	if !s.NewestUpdated.Equal(new_) {
		t.Fatalf("expected newest %v, got %v", new_, s.NewestUpdated)
	}
}
