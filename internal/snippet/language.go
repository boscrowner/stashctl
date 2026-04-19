package snippet

import "strings"

// KnownLanguages is the set of recognized language identifiers.
var KnownLanguages = []string{
	"go", "python", "javascript", "typescript", "rust", "bash",
	"sh", "sql", "yaml", "json", "toml", "markdown", "ruby",
	"java", "c", "cpp", "csharp", "php", "swift", "kotlin",
}

// NormalizeLanguage lowercases and trims a language string.
func NormalizeLanguage(lang string) string {
	return strings.ToLower(strings.TrimSpace(lang))
}

// IsKnownLanguage reports whether lang is in the known list.
func IsKnownLanguage(lang string) bool {
	norm := NormalizeLanguage(lang)
	for _, l := range KnownLanguages {
		if l == norm {
			return true
		}
	}
	return false
}

// SuggestLanguage attempts to infer a language from a file extension.
func SuggestLanguage(filename string) string {
	ext := strings.ToLower(filename)
	if i := strings.LastIndex(ext, "."); i >= 0 {
		ext = ext[i+1:]
	}
	switch ext {
	case "go":
		return "go"
	case "py":
		return "python"
	case "js":
		return "javascript"
	case "ts":
		return "typescript"
	case "rs":
		return "rust"
	case "sh", "bash":
		return "bash"
	case "sql":
		return "sql"
	case "rb":
		return "ruby"
	case "java":
		return "java"
	case "kt":
		return "kotlin"
	case "swift":
		return "swift"
	case "yaml", "yml":
		return "yaml"
	case "json":
		return "json"
	case "toml":
		return "toml"
	case "md":
		return "markdown"
	}
	return ""
}
