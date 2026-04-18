package export

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/stashctl/internal/snippet"
)

// Format represents the output format for exported snippets.
type Format string

const (
	FormatJSON     Format = "json"
	FormatMarkdown Format = "markdown"
)

// Snippets writes the given snippets to w in the specified format.
func Snippets(w io.Writer, snippets []*snippet.Snippet, format Format) error {
	switch format {
	case FormatJSON:
		return exportJSON(w, snippets)
	case FormatMarkdown:
		return exportMarkdown(w, snippets)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

func exportJSON(w io.Writer, snippets []*snippet.Snippet) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(snippets)
}

func exportMarkdown(w io.Writer, snippets []*snippet.Snippet) error {
	for i, s := range snippets {
		if i > 0 {
			fmt.Fprintln(w)
		}
		fmt.Fprintf(w, "## %s\n", s.Title)
		if len(s.Tags) > 0 {
			fmt.Fprintf(w, "**Tags:** %s\n\n", strings.Join(s.Tags, ", "))
		}
		if s.Language != "" {
			fmt.Fprintf(w, "```%s\n%s\n```\n", s.Language, s.Body)
		} else {
			fmt.Fprintf(w, "```\n%s\n```\n", s.Body)
		}
	}
	return nil
}
