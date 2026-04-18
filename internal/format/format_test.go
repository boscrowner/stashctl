package format_test

import (
	"strings"
	"testing"

	"github.com/user/stashctl/internal/format"
	"github.com/user/stashctl/internal/snippet"
)

func makeSnippet() snippet.Snippet {
	return snippet.Snippet{
		ID:    "abc123",
		Title: "Hello World",
		Body:  "fmt.Println(\"hello\")",
		Tags:  []string{"go", "example"},
	}
}

func TestSnippetSummaryNoColor(t *testing.T) {
	s := makeSnippet()
	out := format.SnippetSummary(s, false)
	if !strings.Contains(out, s.ID) {
		t.Errorf("expected ID in summary, got: %s", out)
	}
	if !strings.Contains(out, s.Title) {
		t.Errorf("expected title in summary, got: %s", out)
	}
	if !strings.Contains(out, "go") || !strings.Contains(out, "example") {
		t.Errorf("expected tags in summary, got: %s", out)
	}
}

func TestSnippetSummaryNoTags(t *testing.T) {
	s := makeSnippet()
	s.Tags = nil
	out := format.SnippetSummary(s, false)
	if strings.Contains(out, "[") {
		t.Errorf("expected no tag brackets when tags empty, got: %s", out)
	}
}

func TestSnippetDetailNoColor(t *testing.T) {
	s := makeSnippet()
	out := format.SnippetDetail(s, false)
	for _, want := range []string{s.ID, s.Title, s.Body, "go", "example"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in detail output, got:\n%s", want, out)
		}
	}
}

func TestSnippetDetailColorized(t *testing.T) {
	s := makeSnippet()
	out := format.SnippetDetail(s, true)
	if !strings.Contains(out, format.ColorGray) {
		t.Errorf("expected color codes in colorized output")
	}
	if !strings.Contains(out, s.Title) {
		t.Errorf("expected title in colorized detail")
	}
}
