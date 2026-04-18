package export

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/stashctl/internal/snippet"
)

func makeSnippet(id, title, body, lang string, tags []string) *snippet.Snippet {
	return &snippet.Snippet{
		ID:        id,
		Title:     title,
		Body:      body,
		Language:  lang,
		Tags:      tags,
		CreatedAt: time.Now(),
	}
}

func TestExportJSON(t *testing.T) {
	snippets := []*snippet.Snippet{
		makeSnippet("1", "Hello", "fmt.Println()", "go", []string{"go", "print"}),
	}
	var buf bytes.Buffer
	if err := Snippets(&buf, snippets, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []*snippet.Snippet
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(out) != 1 || out[0].Title != "Hello" {
		t.Errorf("unexpected output: %+v", out)
	}
}

func TestExportMarkdown(t *testing.T) {
	snippets := []*snippet.Snippet{
		makeSnippet("1", "Hello", "fmt.Println()", "go", []string{"go"}),
	}
	var buf bytes.Buffer
	if err := Snippets(&buf, snippets, FormatMarkdown); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "## Hello") {
		t.Errorf("expected title in markdown, got: %s", out)
	}
	if !strings.Contains(out, "```go") {
		t.Errorf("expected go code block, got: %s", out)
	}
}

func TestExportUnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Snippets(&buf, nil, Format("xml"))
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExportMarkdownMultiple(t *testing.T) {
	snippets := []*snippet.Snippet{
		makeSnippet("1", "First", "body1", "", nil),
		makeSnippet("2", "Second", "body2", "bash", []string{"shell"}),
	}
	var buf bytes.Buffer
	if err := Snippets(&buf, snippets, FormatMarkdown); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "## First") || !strings.Contains(out, "## Second") {
		t.Errorf("expected both titles in output: %s", out)
	}
}
