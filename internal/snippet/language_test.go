package snippet_test

import (
	"testing"

	"github.com/user/stashctl/internal/snippet"
)

func TestNormalizeLanguage(t *testing.T) {
	cases := []struct {
		input, want string
	}{
		{"Go", "go"},
		{" Python ", "python"},
		{"RUST", "rust"},
		{"", ""},
	}
	for _, c := range cases {
		got := snippet.NormalizeLanguage(c.input)
		if got != c.want {
			t.Errorf("NormalizeLanguage(%q) = %q; want %q", c.input, got, c.want)
		}
	}
}

func TestIsKnownLanguage(t *testing.T) {
	if !snippet.IsKnownLanguage("go") {
		t.Error("expected go to be known")
	}
	if !snippet.IsKnownLanguage("Python") {
		t.Error("expected Python (case-insensitive) to be known")
	}
	if snippet.IsKnownLanguage("brainfuck") {
		t.Error("expected brainfuck to be unknown")
	}
	if snippet.IsKnownLanguage("") {
		t.Error("expected empty string to be unknown")
	}
}

func TestSuggestLanguage(t *testing.T) {
	cases := []struct {
		filename, want string
	}{
		{"main.go", "go"},
		{"script.py", "python"},
		{"app.js", "javascript"},
		{"index.ts", "typescript"},
		{"run.sh", "bash"},
		{"query.sql", "sql"},
		{"config.yaml", "yaml"},
		{"config.yml", "yaml"},
		{"data.json", "json"},
		{"README.md", "markdown"},
		{"unknown.xyz", ""},
		{"noextension", ""},
	}
	for _, c := range cases {
		got := snippet.SuggestLanguage(c.filename)
		if got != c.want {
			t.Errorf("SuggestLanguage(%q) = %q; want %q", c.filename, got, c.want)
		}
	}
}
