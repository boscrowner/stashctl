package format

import (
	"fmt"
	"strings"

	"github.com/user/stashctl/internal/snippet"
)

const (
	ColorReset  = "\033[0m"
	ColorCyan   = "\033[36m"
	ColorYellow = "\033[33m"
	ColorGreen  = "\033[32m"
	ColorGray   = "\033[90m"
)

// SnippetSummary returns a compact single-line representation of a snippet.
func SnippetSummary(s snippet.Snippet, colorize bool) string {
	tags := ""
	if len(s.Tags) > 0 {
		tags = fmt.Sprintf(" [%s]", strings.Join(s.Tags, ", "))
	}
	if colorize {
		return fmt.Sprintf("%s%s%s  %s%s%s%s%s",
			ColorCyan, s.ID, ColorReset,
			ColorGreen, s.Title, ColorReset,
			ColorYellow, tags+ColorReset,
		)
	}
	return fmt.Sprintf("%s  %s%s", s.ID, s.Title, tags)
}

// SnippetDetail returns a full multi-line representation of a snippet.
func SnippetDetail(s snippet.Snippet, colorize bool) string {
	var sb strings.Builder
	if colorize {
		sb.WriteString(fmt.Sprintf("%sID:%s    %s\n", ColorGray, ColorReset, s.ID))
		sb.WriteString(fmt.Sprintf("%sTitle:%s %s\n", ColorGray, ColorReset, s.Title))
		if len(s.Tags) > 0 {
			sb.WriteString(fmt.Sprintf("%sTags:%s  %s\n", ColorGray, ColorReset, strings.Join(s.Tags, ", ")))
		}
		sb.WriteString(fmt.Sprintf("%sBody:%s\n%s\n", ColorGray, ColorReset, s.Body))
	} else {
		sb.WriteString(fmt.Sprintf("ID:    %s\n", s.ID))
		sb.WriteString(fmt.Sprintf("Title: %s\n", s.Title))
		if len(s.Tags) > 0 {
			sb.WriteString(fmt.Sprintf("Tags:  %s\n", strings.Join(s.Tags, ", ")))
		}
		sb.WriteString(fmt.Sprintf("Body:\n%s\n", s.Body))
	}
	return sb.String()
}
