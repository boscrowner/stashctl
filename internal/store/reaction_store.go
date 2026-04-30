package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/user/stashctl/internal/snippet"
)

const reactionsFile = "reactions.json"

func (s *Store) reactionsPath() string {
	return filepath.Join(s.dir, reactionsFile)
}

func (s *Store) LoadReactions() ([]snippet.Reaction, error) {
	data, err := os.ReadFile(s.reactionsPath())
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var reactions []snippet.Reaction
	if err := json.Unmarshal(data, &reactions); err != nil {
		return nil, err
	}
	return reactions, nil
}

func (s *Store) SaveReactions(reactions []snippet.Reaction) error {
	data, err := json.MarshalIndent(reactions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.reactionsPath(), data, 0o644)
}

func (s *Store) AddReaction(r snippet.Reaction) error {
	reactions, err := s.LoadReactions()
	if err != nil {
		return err
	}
	reactions = append(reactions, r)
	return s.SaveReactions(reactions)
}

func (s *Store) DeleteReaction(id string) error {
	reactions, err := s.LoadReactions()
	if err != nil {
		return err
	}
	updated, err := snippet.RemoveReaction(reactions, id)
	if err != nil {
		return err
	}
	return s.SaveReactions(updated)
}

func (s *Store) ListReactions(snippetID string) ([]snippet.Reaction, error) {
	reactions, err := s.LoadReactions()
	if err != nil {
		return nil, err
	}
	return snippet.ReactionsFor(reactions, snippetID), nil
}
