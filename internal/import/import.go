// Package importer provides functionality to import snippets from various formats.
package importer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/stashctl/internal/snippet"
)

// Result holds the outcome of an import operation.
type Result struct {
	Imported int
	Skipped  int
	Errors   []string
}

// FromFile detects the format by file extension and imports snippets.
func FromFile(path string) ([]*snippet.Snippet, Result, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return fromJSON(path)
	case ".md", ".markdown":
		return nil, Result{}, fmt.Errorf("markdown import not yet supported")
	default:
		return nil, Result{}, fmt.Errorf("unsupported import format: %s", ext)
	}
}

// jsonSnippet mirrors the JSON export shape.
type jsonSnippet struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Language string   `json:"language"`
	Tags     []string `json:"tags"`
}

func fromJSON(path string) ([]*snippet.Snippet, Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, Result{}, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var raw []jsonSnippet
	if err := json.NewDecoder(f).Decode(&raw); err != nil {
		return nil, Result{}, fmt.Errorf("decode json: %w", err)
	}

	var snippets []*snippet.Snippet
	var res Result
	for _, r := range raw {
		if r.Title == "" || r.Content == "" {
			res.Skipped++
			res.Errors = append(res.Errors, fmt.Sprintf("skipped entry with empty title or content"))
			continue
		}
		s := snippet.New(r.Title, r.Content, r.Language, r.Tags)
		snippets = append(snippets, s)
		res.Imported++
	}
	return snippets, res, nil
}
