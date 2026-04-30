package snippet_test

import (
	"strings"
	"testing"

	"github.com/user/stashctl/internal/snippet"
)

func makeChecklistItem(t *testing.T, snippetID, text string) snippet.ChecklistItem {
	t.Helper()
	item, err := snippet.NewChecklistItem(snippetID, text)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return item
}

func TestNewChecklistItemValid(t *testing.T) {
	item := makeChecklistItem(t, "snip-1", "Write tests")
	if item.ID == "" {
		t.Error("expected non-empty ID")
	}
	if item.Done {
		t.Error("new item should not be done")
	}
	if item.Text != "Write tests" {
		t.Errorf("unexpected text: %s", item.Text)
	}
}

func TestNewChecklistItemEmptySnippetID(t *testing.T) {
	_, err := snippet.NewChecklistItem("", "Do something")
	if err == nil {
		t.Error("expected error for empty snippet ID")
	}
}

func TestNewChecklistItemEmptyText(t *testing.T) {
	_, err := snippet.NewChecklistItem("snip-1", "")
	if err == nil {
		t.Error("expected error for empty text")
	}
}

func TestNewChecklistItemTextTooLong(t *testing.T) {
	_, err := snippet.NewChecklistItem("snip-1", strings.Repeat("x", 201))
	if err == nil {
		t.Error("expected error for text exceeding 200 chars")
	}
}

func TestChecklistFor(t *testing.T) {
	items := []snippet.ChecklistItem{
		makeChecklistItem(t, "snip-1", "Task A"),
		makeChecklistItem(t, "snip-2", "Task B"),
		makeChecklistItem(t, "snip-1", "Task C"),
	}
	got := snippet.ChecklistFor(items, "snip-1")
	if len(got) != 2 {
		t.Errorf("expected 2 items, got %d", len(got))
	}
}

func TestCompleteItem(t *testing.T) {
	items := []snippet.ChecklistItem{makeChecklistItem(t, "snip-1", "Task A")}
	updated, err := snippet.CompleteItem(items, items[0].ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !updated[0].Done {
		t.Error("expected item to be marked done")
	}
}

func TestCompleteItemNotFound(t *testing.T) {
	items := []snippet.ChecklistItem{makeChecklistItem(t, "snip-1", "Task A")}
	_, err := snippet.CompleteItem(items, "nonexistent")
	if err == nil {
		t.Error("expected error for missing item")
	}
}

func TestRemoveChecklistItem(t *testing.T) {
	item := makeChecklistItem(t, "snip-1", "Task A")
	items := []snippet.ChecklistItem{item}
	updated, err := snippet.RemoveChecklistItem(items, item.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updated) != 0 {
		t.Errorf("expected 0 items after removal, got %d", len(updated))
	}
}

func TestPendingItems(t *testing.T) {
	items := []snippet.ChecklistItem{
		makeChecklistItem(t, "snip-1", "Task A"),
		makeChecklistItem(t, "snip-1", "Task B"),
	}
	items, _ = snippet.CompleteItem(items, items[0].ID)
	pending := snippet.PendingItems(items, "snip-1")
	if len(pending) != 1 {
		t.Errorf("expected 1 pending item, got %d", len(pending))
	}
	if pending[0].Text != "Task B" {
		t.Errorf("unexpected pending item: %s", pending[0].Text)
	}
}
