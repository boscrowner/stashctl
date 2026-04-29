package snippet

import (
	"testing"
)

func makeWebhook(t *testing.T, snippetID, rawURL string, events []WebhookEvent) Webhook {
	t.Helper()
	w, err := NewWebhook(snippetID, rawURL, events, "")
	if err != nil {
		t.Fatalf("makeWebhook: %v", err)
	}
	return w
}

func TestNewWebhookValid(t *testing.T) {
	w, err := NewWebhook("s1", "https://example.com/hook", []WebhookEvent{WebhookEventCreated}, "mysecret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.ID == "" {
		t.Error("expected non-empty ID")
	}
	if w.Secret != "mysecret" {
		t.Errorf("expected secret mysecret, got %s", w.Secret)
	}
}

func TestNewWebhookEmptySnippetID(t *testing.T) {
	_, err := NewWebhook("", "https://example.com/hook", []WebhookEvent{WebhookEventCreated}, "")
	if err == nil {
		t.Error("expected error for empty snippet_id")
	}
}

func TestNewWebhookInvalidURL(t *testing.T) {
	_, err := NewWebhook("s1", "not-a-url", []WebhookEvent{WebhookEventCreated}, "")
	if err == nil {
		t.Error("expected error for invalid URL")
	}
}

func TestNewWebhookNonHTTPURL(t *testing.T) {
	_, err := NewWebhook("s1", "ftp://example.com/hook", []WebhookEvent{WebhookEventCreated}, "")
	if err == nil {
		t.Error("expected error for non-http URL")
	}
}

func TestNewWebhookNoEvents(t *testing.T) {
	_, err := NewWebhook("s1", "https://example.com/hook", nil, "")
	if err == nil {
		t.Error("expected error for empty events")
	}
}

func TestNewWebhookUnknownEvent(t *testing.T) {
	_, err := NewWebhook("s1", "https://example.com/hook", []WebhookEvent{"archived"}, "")
	if err == nil {
		t.Error("expected error for unknown event")
	}
}

func TestWebhooksFor(t *testing.T) {
	w1 := makeWebhook(t, "s1", "https://a.com", []WebhookEvent{WebhookEventCreated})
	w2 := makeWebhook(t, "s2", "https://b.com", []WebhookEvent{WebhookEventDeleted})
	w3 := makeWebhook(t, "s1", "https://c.com", []WebhookEvent{WebhookEventUpdated})
	all := []Webhook{w1, w2, w3}
	got := WebhooksFor(all, "s1")
	if len(got) != 2 {
		t.Fatalf("expected 2 webhooks for s1, got %d", len(got))
	}
}

func TestRemoveWebhook(t *testing.T) {
	w1 := makeWebhook(t, "s1", "https://a.com", []WebhookEvent{WebhookEventCreated})
	w2 := makeWebhook(t, "s1", "https://b.com", []WebhookEvent{WebhookEventDeleted})
	all := []Webhook{w1, w2}
	updated, err := RemoveWebhook(all, w1.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updated) != 1 {
		t.Fatalf("expected 1 webhook after removal, got %d", len(updated))
	}
}

func TestRemoveWebhookNotFound(t *testing.T) {
	w := makeWebhook(t, "s1", "https://a.com", []WebhookEvent{WebhookEventCreated})
	_, err := RemoveWebhook([]Webhook{w}, "nonexistent")
	if err == nil {
		t.Error("expected error for missing webhook")
	}
}
