package snippet

import (
	"testing"
)

func TestNewTemplateValid(t *testing.T) {
	tmpl, err := NewTemplate("http handler", "go", "func {{Name}}(w http.ResponseWriter, r *http.Request) {\n}", []string{"http", "go"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tmpl.Name != "http handler" {
		t.Errorf("expected name 'http handler', got %q", tmpl.Name)
	}
	if tmpl.Language != "go" {
		t.Errorf("expected language 'go', got %q", tmpl.Language)
	}
	if tmpl.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestNewTemplateEmptyName(t *testing.T) {
	_, err := NewTemplate("", "go", "some body", nil)
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestNewTemplateEmptyBody(t *testing.T) {
	_, err := NewTemplate("tmpl", "go", "   ", nil)
	if err == nil {
		t.Error("expected error for empty body")
	}
}

func TestInstantiateSubstitution(t *testing.T) {
	tmpl, _ := NewTemplate("greet", "go", "fmt.Println(\"Hello, {{Name}}!\")", []string{"fmt"})
	s, err := tmpl.Instantiate("greet john", map[string]string{"Name": "John"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Content != "fmt.Println(\"Hello, John!\")" {
		t.Errorf("unexpected content: %q", s.Content)
	}
	if s.Title != "greet john" {
		t.Errorf("unexpected title: %q", s.Title)
	}
}

func TestInstantiateEmptyTitle(t *testing.T) {
	tmpl, _ := NewTemplate("greet", "go", "body", nil)
	_, err := tmpl.Instantiate("", nil)
	if err == nil {
		t.Error("expected error for empty title")
	}
}

func TestInstantiateTagsInherited(t *testing.T) {
	tmpl, _ := NewTemplate("base", "python", "print({{msg}})", []string{"print", "python"})
	s, _ := tmpl.Instantiate("hello", map[string]string{"msg": "'hi'"})
	if len(s.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(s.Tags))
	}
}
