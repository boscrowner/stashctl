package format_test

import (
	"strings"
	"testing"

	"github.com/user/stashctl/internal/format"
)

func TestLanguageListContainsGo(t *testing.T) {
	out := format.LanguageList()
	if !strings.Contains(out, "go") {
		t.Error("expected language list to contain 'go'")
	}
	if !strings.Contains(out, "Supported languages:") {
		t.Error("expected header line")
	}
}

func TestLanguageBadgePlain(t *testing.T) {
	badge := format.LanguageBadge("python", false)
	if badge != "[python]" {
		t.Errorf("got %q; want [python]", badge)
	}
}

func TestLanguageBadgeEmpty(t *testing.T) {
	badge := format.LanguageBadge("", false)
	if badge != "[plain]" {
		t.Errorf("got %q; want [plain]", badge)
	}
}

func TestLanguageBadgeColor(t *testing.T) {
	badge := format.LanguageBadge("go", true)
	if !strings.Contains(badge, "go") {
		t.Error("colored badge should still contain language name")
	}
	if !strings.Contains(badge, "\033[") {
		t.Error("colored badge should contain ANSI escape")
	}
}
