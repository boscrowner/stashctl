package snippet

import (
	"testing"
)

func TestNormalizeTagsDeduplicates(t *testing.T) {
	result := NormalizeTags([]string{"Go", "go", "GO"})
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
	result := NormalizeTags([]string{"", "  ", "go"})
	if len(result) != 1 || result[0] != "go" {
		t.Fatalf("expected [go], got %v", result)
	}
}

func TestParseTags(t *testing.T) {
	result := ParseTags("Go,CLI, tools ,GO")
	expected := []string{"cli", "go", "tools"}
	if !TagsEqual(result, expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestTagsEqual(t *testing.T) {
	if !TagsEqual([]string{"b", "a"}, []string{"A", "B"}) {
		t.Fatal("expected equal")
	}
	if TagsEqual([]string{"a"}, []string{"a", "b"}) {
		t.Fatal("expected not equal")
	}
}
