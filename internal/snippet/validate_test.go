package snippet

import (
	"strings"
	"testing"
)

func makeValidSnippet() *Snippet {
	s, _ := New("Test Title", "fmt.Println()", "go", []string{"demo"})
	return s
}

func TestValidateOK(t *testing.T) {
	s := makeValidSnippet()
	if err := Validate(s); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateEmptyTitle(t *testing.T) {
	s := makeValidSnippet()
	s.Title = "   "
	if err := Validate(s); err != ErrEmptyTitle {
		t.Fatalf("expected ErrEmptyTitle, got %v", err)
	}
}

func TestValidateTitleTooLong(t *testing.T) {
	s := makeValidSnippet()
	s.Title = strings.Repeat("a", 201)
	if err := Validate(s); err != ErrTitleTooLong {
		t.Fatalf("expected ErrTitleTooLong, got %v", err)
	}
}

func TestValidateEmptyContent(t *testing.T) {
	s := makeValidSnippet()
	s.Content = ""
	if err := Validate(s); err != ErrEmptyContent {
		t.Fatalf("expected ErrEmptyContent, got %v", err)
	}
}

func TestValidateInvalidTagWithSpace(t *testing.T) {
	s := makeValidSnippet()
	s.Tags = []string{"bad tag"}
	if err := Validate(s); err != ErrInvalidTag {
		t.Fatalf("expected ErrInvalidTag, got %v", err)
	}
}

func TestValidateInvalidTagWithComma(t *testing.T) {
	s := makeValidSnippet()
	s.Tags = []string{"bad,tag"}
	if err := Validate(s); err != ErrInvalidTag {
		t.Fatalf("expected ErrInvalidTag, got %v", err)
	}
}
