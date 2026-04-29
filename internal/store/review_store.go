package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/user/stashctl/internal/snippet"
)

const reviewFile = "reviews.json"

func (s *Store) reviewPath() string {
	return filepath.Join(s.dir, reviewFile)
}

func (s *Store) LoadReviews() ([]snippet.Review, error) {
	data, err := os.ReadFile(s.reviewPath())
	if errors.Is(err, os.ErrNotExist) {
		return []snippet.Review{}, nil
	}
	if err != nil {
		return nil, err
	}
	var reviews []snippet.Review
	if err := json.Unmarshal(data, &reviews); err != nil {
		return nil, err
	}
	return reviews, nil
}

func (s *Store) SaveReviews(reviews []snippet.Review) error {
	data, err := json.MarshalIndent(reviews, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.reviewPath(), data, 0o644)
}

func (s *Store) AddReview(r snippet.Review) error {
	reviews, err := s.LoadReviews()
	if err != nil {
		return err
	}
	reviews = append(reviews, r)
	return s.SaveReviews(reviews)
}

func (s *Store) DeleteReview(id string) error {
	reviews, err := s.LoadReviews()
	if err != nil {
		return err
	}
	updated, err := snippet.RemoveReview(reviews, id)
	if err != nil {
		return err
	}
	return s.SaveReviews(updated)
}

func (s *Store) ReviewsForSnippet(snippetID string) ([]snippet.Review, error) {
	reviews, err := s.LoadReviews()
	if err != nil {
		return nil, err
	}
	return snippet.ReviewsFor(reviews, snippetID), nil
}
