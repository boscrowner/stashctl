package snippet

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// SubscriptionEvent represents the type of event a subscription listens for.
type SubscriptionEvent string

const (
	EventUpdated  SubscriptionEvent = "updated"
	EventDeleted  SubscriptionEvent = "deleted"
	EventTagged   SubscriptionEvent = "tagged"
	EventArchived SubscriptionEvent = "archived"
)

var knownSubscriptionEvents = map[SubscriptionEvent]bool{
	EventUpdated:  true,
	EventDeleted:  true,
	EventTagged:   true,
	EventArchived: true,
}

// Subscription represents a named subscription to events on a snippet.
type Subscription struct {
	ID         string            `json:"id"`
	SnippetID  string            `json:"snippet_id"`
	Subscriber string            `json:"subscriber"`
	Event      SubscriptionEvent `json:"event"`
	Note       string            `json:"note,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
}

// NewSubscription creates a validated Subscription.
func NewSubscription(snippetID, subscriber string, event SubscriptionEvent, note string) (Subscription, error) {
	if snippetID == "" {
		return Subscription{}, errors.New("subscription: snippet_id is required")
	}
	if subscriber == "" {
		return Subscription{}, errors.New("subscription: subscriber is required")
	}
	if !knownSubscriptionEvents[event] {
		return Subscription{}, errors.New("subscription: unknown event type")
	}
	if len(note) > 200 {
		return Subscription{}, errors.New("subscription: note exceeds 200 characters")
	}
	return Subscription{
		ID:         uuid.NewString(),
		SnippetID:  snippetID,
		Subscriber: subscriber,
		Event:      event,
		Note:       note,
		CreatedAt:  time.Now().UTC(),
	}, nil
}

// SubscriptionsFor returns all subscriptions for a given snippet ID.
func SubscriptionsFor(subs []Subscription, snippetID string) []Subscription {
	var out []Subscription
	for _, s := range subs {
		if s.SnippetID == snippetID {
			out = append(out, s)
		}
	}
	return out
}

// RemoveSubscription removes a subscription by ID.
func RemoveSubscription(subs []Subscription, id string) ([]Subscription, bool) {
	for i, s := range subs {
		if s.ID == id {
			return append(subs[:i], subs[i+1:]...), true
		}
	}
	return subs, false
}

// FindSubscription returns a subscription by ID.
func FindSubscription(subs []Subscription, id string) (Subscription, bool) {
	for _, s := range subs {
		if s.ID == id {
			return s, true
		}
	}
	return Subscription{}, false
}
