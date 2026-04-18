package snippet

import (
	"testing"
)

func TestNormalizeTagsDeduplicates(t *testing.T) {
	tags := NormalizeTags([]string{"go", "Go", "GO"})
	if len(tags) != 1 || tags[0] != "go" {
		t.Fatalf("expected [go], got %v", tags)
	}
}

func TestNormalizeTagsSorted(t *testing.T) {
	tags := NormalizeTags([]string{"zebra", "alpha", "mango"})
	if tags[0] != "alpha" || tags[1] != "mango" || tags[2] != "zebra" {
		t.Fatalf("expected sorted tags, got %v", tags)
	}
}

func TestNormalizeTagsStripsEmpty(t *testing.T) {
	tags := NormalizeTags([]string{"go", "", "  ", "cli"})
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %v", tags)
	}
}

func TestParseTags(t *testing.T) {
	tags := ParseTags("Go, CLI, go")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %v", tags)
	}
	if tags[0] != "cli" || tags[1] != "go" {
		t.Fatalf("unexpected tags: %v", tags)
	}
}

func TestTagsEqual(t *testing.T) {
	if !TagsEqual([]string{"go", "cli"}, []string{"CLI", "Go"}) {
		t.Fatal("expected tags to be equal")
	}
	if TagsEqual([]string{"go"}, []string{"go", "cli"}) {
		t.Fatal("expected tags to be unequal")
	}
}
