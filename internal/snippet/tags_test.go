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
	tags := NormalizeTags([]string{"zebra", "apple", "mango"})
	expected := []string{"apple", "mango", "zebra"}
	if !TagsEqual(tags, expected) {
		t.Fatalf("expected %v, got %v", expected, tags)
	}
}

func TestNormalizeTagsStripsEmpty(t *testing.T) {
	tags := NormalizeTags([]string{"go", "", "  ", "cli"})
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %v", tags)
	}
}

func TestParseTags(t *testing.T) {
	tags := ParseTags("go, CLI , Go")
	expected := []string{"cli", "go"}
	if !TagsEqual(tags, expected) {
		t.Fatalf("expected %v, got %v", expected, tags)
	}
}

func TestTagsEqual(t *testing.T) {
	if !TagsEqual([]string{"b", "a"}, []string{"a", "b"}) {
		t.Fatal("expected equal")
	}
	if TagsEqual([]string{"a"}, []string{"a", "b"}) {
		t.Fatal("expected not equal")
	}
}
