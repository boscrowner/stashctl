package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/user/stashctl/internal/snippet"
)

const subscriptionFile = "subscriptions.json"

func (s *Store) subscriptionPath() string {
	return filepath.Join(s.dir, subscriptionFile)
}

func (s *Store) loadSubscriptions() ([]snippet.Subscription, error) {
	data, err := os.ReadFile(s.subscriptionPath())
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var subs []snippet.Subscription
	if err := json.Unmarshal(data, &subs); err != nil {
		return nil, err
	}
	return subs, nil
}

func (s *Store) saveSubscriptions(subs []snippet.Subscription) error {
	data, err := json.MarshalIndent(subs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.subscriptionPath(), data, 0o644)
}

// AddSubscription persists a new subscription.
func (s *Store) AddSubscription(sub snippet.Subscription) error {
	subs, err := s.loadSubscriptions()
	if err != nil {
		return err
	}
	subs = append(subs, sub)
	return s.saveSubscriptions(subs)
}

// ListSubscriptions returns all subscriptions for a snippet.
func (s *Store) ListSubscriptions(snippetID string) ([]snippet.Subscription, error) {
	subs, err := s.loadSubscriptions()
	if err != nil {
		return nil, err
	}
	return snippet.SubscriptionsFor(subs, snippetID), nil
}

// DeleteSubscription removes a subscription by ID.
func (s *Store) DeleteSubscription(id string) error {
	subs, err := s.loadSubscriptions()
	if err != nil {
		return err
	}
	updated, ok := snippet.RemoveSubscription(subs, id)
	if !ok {
		return errors.New("subscription not found")
	}
	return s.saveSubscriptions(updated)
}

// AllSubscriptions returns every stored subscription.
func (s *Store) AllSubscriptions() ([]snippet.Subscription, error) {
	return s.loadSubscriptions()
}
