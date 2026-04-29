package snippet

import (
	"errors"
	"net/url"
	"time"

	"github.com/google/uuid"
)

// WebhookEvent represents the kind of event that triggers a webhook.
type WebhookEvent string

const (
	WebhookEventCreated WebhookEvent = "created"
	WebhookEventUpdated WebhookEvent = "updated"
	WebhookEventDeleted WebhookEvent = "deleted"
)

// Webhook represents a registered HTTP callback for snippet lifecycle events.
type Webhook struct {
	ID        string       `json:"id"`
	SnippetID string       `json:"snippet_id"`
	URL       string       `json:"url"`
	Events    []WebhookEvent `json:"events"`
	Secret    string       `json:"secret,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
}

// NewWebhook creates and validates a new Webhook.
func NewWebhook(snippetID, rawURL string, events []WebhookEvent, secret string) (Webhook, error) {
	if snippetID == "" {
		return Webhook{}, errors.New("snippet_id is required")
	}
	if rawURL == "" {
		return Webhook{}, errors.New("url is required")
	}
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return Webhook{}, errors.New("url must be a valid http or https URL")
	}
	if len(events) == 0 {
		return Webhook{}, errors.New("at least one event is required")
	}
	for _, e := range events {
		if !isKnownWebhookEvent(e) {
			return Webhook{}, errors.New("unknown event: " + string(e))
		}
	}
	return Webhook{
		ID:        uuid.NewString(),
		SnippetID: snippetID,
		URL:       rawURL,
		Events:    events,
		Secret:    secret,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// WebhooksFor returns all webhooks registered for a given snippet.
func WebhooksFor(all []Webhook, snippetID string) []Webhook {
	var out []Webhook
	for _, w := range all {
		if w.SnippetID == snippetID {
			out = append(out, w)
		}
	}
	return out
}

// RemoveWebhook removes a webhook by ID, returning the updated slice.
func RemoveWebhook(all []Webhook, id string) ([]Webhook, error) {
	for i, w := range all {
		if w.ID == id {
			return append(all[:i], all[i+1:]...), nil
		}
	}
	return nil, errors.New("webhook not found: " + id)
}

func isKnownWebhookEvent(e WebhookEvent) bool {
	switch e {
	case WebhookEventCreated, WebhookEventUpdated, WebhookEventDeleted:
		return true
	}
	return false
}
