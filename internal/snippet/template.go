package snippet

import (
	"errors"
	"strings"
	"time"
)

// Template represents a reusable snippet scaffold.
type Template struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Language string    `json:"language"`
	Tags     []string  `json:"tags"`
	Body     string    `json:"body"`
	Created  time.Time `json:"created"`
}

// NewTemplate creates a new Template with the given name, language, tags and body.
func NewTemplate(name, language, body string, tags []string) (*Template, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("template name must not be empty")
	}
	if strings.TrimSpace(body) == "" {
		return nil, errors.New("template body must not be empty")
	}
	return &Template{
		ID:       generateID(),
		Name:     name,
		Language: NormalizeLanguage(language),
		Tags:     NormalizeTags(tags),
		Body:     body,
		Created:  time.Now().UTC(),
	}, nil
}

// Instantiate creates a Snippet from the template, substituting placeholders
// in the body with the provided values map ({{KEY}} -> value).
func (t *Template) Instantiate(title string, values map[string]string) (*Snippet, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("snippet title must not be empty")
	}
	body := t.Body
	for k, v := range values {
		body = strings.ReplaceAll(body, "{{"+k+"}}", v)
	}
	s := New(title, body, t.Language, append([]string(nil), t.Tags...))
	return s, nil
}
