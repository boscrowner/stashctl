package snippet

import (
	"testing"
)

func TestNormalizeTagsDeduplicates(t *testing.T) {
	result := NormalizeTags([]string{"go", "Go", "GO"})
	if len(result) != 1 || result[0] != "go" {
		t.Fatalf("expected [go], got %v", result)
	}
}

func TestNormalizeTagsSorted(t *testing.T) {
	result := NormalizeTags([]string{"zebra", "apple", "mango"})
	expected := []string{"apple", "mango", "zebra"}
	if !TagsEqual(result, expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestNormalizeTagsStripsEmpty(t *testing.T) {
	result := NormalizeTags([]string{"go", "", "  ", "cli"})
	if len(result) != 2 {
		t.Fatalf("expected 2 tags, got %v", result)
	}
}

func TestParseTags(t *testing.T) {
	result := ParseTags("Go, CLI,  search,GO")
	expected := []string{"cli", "go", "search"}
	if !TagsEqual(result, expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestTagsEqual(t *testing.T) {
	if !TagsEqual([]string{"a", "b"}, []string{"b", "a"}) {
		t.Fatal("expected tags to be equal regardless of order")
	}
	if TagsEqual([]string{"a"}, []string{"a", "b"}) {
		t.Fatal("expected tags to be unequal")
	}
}
