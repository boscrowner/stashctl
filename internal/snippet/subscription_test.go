package snippet_test

import (
	"strings"
	"testing"

	"github.com/user/stashctl/internal/snippet"
)

func makeSub(snippetID, subscriber string, event snippet.SubscriptionEvent) snippet.Subscription {
	s, err := snippet.NewSubscription(snippetID, subscriber, event, "")
	if err != nil {
		panic(err)
	}
	return s
}

func TestNewSubscriptionValid(t *testing.T) {
	s, err := snippet.NewSubscription("s1", "alice", snippet.EventUpdated, "notify me")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.ID == "" {
		t.Error("expected non-empty ID")
	}
	if s.SnippetID != "s1" || s.Subscriber != "alice" {
		t.Error("fields not set correctly")
	}
}

func TestNewSubscriptionEmptySnippetID(t *testing.T) {
	_, err := snippet.NewSubscription("", "alice", snippet.EventUpdated, "")
	if err == nil {
		t.Error("expected error for empty snippet_id")
	}
}

func TestNewSubscriptionEmptySubscriber(t *testing.T) {
	_, err := snippet.NewSubscription("s1", "", snippet.EventUpdated, "")
	if err == nil {
		t.Error("expected error for empty subscriber")
	}
}

func TestNewSubscriptionUnknownEvent(t *testing.T) {
	_, err := snippet.NewSubscription("s1", "alice", "unknown", "")
	if err == nil {
		t.Error("expected error for unknown event")
	}
}

func TestNewSubscriptionNoteTooLong(t *testing.T) {
	_, err := snippet.NewSubscription("s1", "alice", snippet.EventUpdated, strings.Repeat("x", 201))
	if err == nil {
		t.Error("expected error for note too long")
	}
}

func TestSubscriptionsFor(t *testing.T) {
	subs := []snippet.Subscription{
		makeSub("s1", "alice", snippet.EventUpdated),
		makeSub("s2", "bob", snippet.EventDeleted),
		makeSub("s1", "carol", snippet.EventTagged),
	}
	result := snippet.SubscriptionsFor(subs, "s1")
	if len(result) != 2 {
		t.Fatalf("expected 2 subscriptions, got %d", len(result))
	}
}

func TestRemoveSubscription(t *testing.T) {
	s := makeSub("s1", "alice", snippet.EventUpdated)
	subs := []snippet.Subscription{s}
	updated, ok := snippet.RemoveSubscription(subs, s.ID)
	if !ok {
		t.Error("expected removal to succeed")
	}
	if len(updated) != 0 {
		t.Error("expected empty slice after removal")
	}
}

func TestRemoveSubscriptionNotFound(t *testing.T) {
	s := makeSub("s1", "alice", snippet.EventUpdated)
	subs := []snippet.Subscription{s}
	_, ok := snippet.RemoveSubscription(subs, "nonexistent")
	if ok {
		t.Error("expected removal to fail for unknown ID")
	}
}

func TestFindSubscription(t *testing.T) {
	s := makeSub("s1", "alice", snippet.EventArchived)
	subs := []snippet.Subscription{s}
	found, ok := snippet.FindSubscription(subs, s.ID)
	if !ok || found.ID != s.ID {
		t.Error("expected to find subscription by ID")
	}
}
