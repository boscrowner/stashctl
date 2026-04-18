package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/user/stashctl/internal/snippet"
)

var ErrNotFound = errors.New("snippet not found")

// Store manages persistence of snippets to a JSON file.
type Store struct {
	path string
}

// New creates a Store backed by the given file path.
func New(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	return &Store{path: path}, nil
}

// load reads all snippets from disk.
func (s *Store) load() ([]*snippet.Snippet, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return []*snippet.Snippet{}, nil
	}
	if err != nil {
		return nil, err
	}
	var snippets []*snippet.Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		return nil, err
	}
	return snippets, nil
}

// save writes all snippets to disk.
func (s *Store) save(snippets []*snippet.Snippet) error {
	data, err := json.MarshalIndent(snippets, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Add persists a new snippet.
func (s *Store) Add(sn *snippet.Snippet) error {
	snippets, err := s.load()
	if err != nil {
		return err
	}
	snippets = append(snippets, sn)
	return s.save(snippets)
}

// List returns all stored snippets.
func (s *Store) List() ([]*snippet.Snippet, error) {
	return s.load()
}

// Delete removes a snippet by ID.
func (s *Store) Delete(id string) error {
	snippets, err := s.load()
	if err != nil {
		return err
	}
	filtered := snippets[:0]
	found := false
	for _, sn := range snippets {
		if sn.ID == id {
			found = true
			continue
		}
		filtered = append(filtered, sn)
	}
	if !found {
		return ErrNotFound
	}
	return s.save(filtered)
}

// FilterByTags returns snippets that have all of the given tags.
func (s *Store) FilterByTags(tags []string) ([]*snippet.Snippet, error) {
	snippets, err := s.load()
	if err != nil {
		return nil, err
	}
	var result []*snippet.Snippet
	for _, sn := range snippets {
		if sn.HasAllTags(tags) {
			result = append(result, sn)
		}
	}
	return result, nil
}
