package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/user/stashctl/internal/snippet"
)

const webhookFile = "webhooks.json"

func (s *Store) webhookPath() string {
	return filepath.Join(s.dir, webhookFile)
}

func (s *Store) loadWebhooks() ([]snippet.Webhook, error) {
	data, err := os.ReadFile(s.webhookPath())
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var out []snippet.Webhook
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Store) saveWebhooks(hooks []snippet.Webhook) error {
	data, err := json.MarshalIndent(hooks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.webhookPath(), data, 0o644)
}

// AddWebhook persists a new webhook.
func (s *Store) AddWebhook(w snippet.Webhook) error {
	hooks, err := s.loadWebhooks()
	if err != nil {
		return err
	}
	hooks = append(hooks, w)
	return s.saveWebhooks(hooks)
}

// ListWebhooks returns all webhooks for a snippet.
func (s *Store) ListWebhooks(snippetID string) ([]snippet.Webhook, error) {
	hooks, err := s.loadWebhooks()
	if err != nil {
		return nil, err
	}
	return snippet.WebhooksFor(hooks, snippetID), nil
}

// DeleteWebhook removes a webhook by ID.
func (s *Store) DeleteWebhook(id string) error {
	hooks, err := s.loadWebhooks()
	if err != nil {
		return err
	}
	updated, err := snippet.RemoveWebhook(hooks, id)
	if err != nil {
		return err
	}
	return s.saveWebhooks(updated)
}

// GetWebhook returns a single webhook by ID, or an error if not found.
func (s *Store) GetWebhook(id string) (snippet.Webhook, error) {
	hooks, err := s.loadWebhooks()
	if err != nil {
		return snippet.Webhook{}, err
	}
	for _, h := range hooks {
		if h.ID == id {
			return h, nil
		}
	}
	return snippet.Webhook{}, errors.New("webhook not found: " + id)
}
