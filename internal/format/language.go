package format

import (
	"fmt"
	"strings"

	"github.com/user/stashctl/internal/snippet"
)

// LanguageList returns a formatted list of all known languages.
func LanguageList() string {
	langs := snippet.KnownLanguages
	var sb strings.Builder
	sb.WriteString("Supported languages:\n")
	for i, l := range langs {
		if i > 0 && i%5 == 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("  %-12s", l))
	}
	sb.WriteString("\n")
	return sb.String()
}

// LanguageBadge returns a short colored badge string for a language.
// When color is false it returns a plain bracketed label.
func LanguageBadge(lang string, color bool) string {
	if lang == "" {
		lang = "plain"
	}
	if !color {
		return fmt.Sprintf("[%s]", lang)
	}
	// simple ANSI cyan for language badge
	return fmt.Sprintf("\033[36m[%s]\033[0m", lang)
}
