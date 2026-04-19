package snippet

import "time"

// Snippet represents a stored code snippet.
type Snippet struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Language  string    `json:"language"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// New creates a new Snippet with a generated ID and normalized tags.
func New(title, content, language string, tags []string) *Snippet {
	now := time.Now().UTC()
	return &Snippet{
		ID:        generateID(),
		Title:     title,
		Content:   content,
		Language:  language,
		Tags:      NormalizeTags(tags),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// HasTag reports whether the snippet has the given tag.
func (s *Snippet) HasTag(tag string) bool {
	for _, t := range s.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// HasAllTags reports whether the snippet has all given tags.
func (s *Snippet) HasAllTags(tags []string) bool {
	for _, tag := range tags {
		if !s.HasTag(tag) {
			return false
		}
	}
	return true
}

// Update applies non-empty fields to the snippet and bumps UpdatedAt.
func (s *Snippet) Update(title, content, language string, tags []string) {
	if title != "" {
		s.Title = title
	}
	if content != "" {
		s.Content = content
	}
	if language != "" {
		s.Language = language
	}
	if tags != nil {
		s.Tags = NormalizeTags(tags)
	}
	s.UpdatedAt = time.Now().UTC()
}
