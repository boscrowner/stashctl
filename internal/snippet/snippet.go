package snippet

import (
	"time"
)

// Snippet represents a stored code snippet with metadata.
type Snippet struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Language    string    `json:"language"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// HasTag returns true if the snippet contains the given tag.
func (s *Snippet) HasTag(tag string) bool {
	for _, t := range s.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// HasAllTags returns true if the snippet contains all of the given tags.
func (s *Snippet) HasAllTags(tags []string) bool {
	for _, tag := range tags {
		if !s.HasTag(tag) {
			return false
		}
	}
	return true
}

// New creates a new Snippet with generated ID and timestamps.
func New(title, content, language string, tags []string) *Snippet {
	now := time.Now().UTC()
	return &Snippet{
		ID:        generateID(),
		Title:     title,
		Content:   content,
		Language:  language,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
