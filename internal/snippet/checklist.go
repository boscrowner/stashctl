package snippet

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ChecklistItem represents a single item in a snippet checklist.
type ChecklistItem struct {
	ID        string    `json:"id"`
	SnippetID string    `json:"snippet_id"`
	Text      string    `json:"text"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewChecklistItem creates a new checklist item for the given snippet.
func NewChecklistItem(snippetID, text string) (ChecklistItem, error) {
	if snippetID == "" {
		return ChecklistItem{}, errors.New("snippet ID must not be empty")
	}
	if text == "" {
		return ChecklistItem{}, errors.New("checklist item text must not be empty")
	}
	if len(text) > 200 {
		return ChecklistItem{}, errors.New("checklist item text must not exceed 200 characters")
	}
	now := time.Now().UTC()
	return ChecklistItem{
		ID:        uuid.NewString(),
		SnippetID: snippetID,
		Text:      text,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// ChecklistFor returns all checklist items belonging to the given snippet.
func ChecklistFor(items []ChecklistItem, snippetID string) []ChecklistItem {
	var result []ChecklistItem
	for _, item := range items {
		if item.SnippetID == snippetID {
			result = append(result, item)
		}
	}
	return result
}

// CompleteItem marks the checklist item with the given ID as done.
func CompleteItem(items []ChecklistItem, id string) ([]ChecklistItem, error) {
	for i, item := range items {
		if item.ID == id {
			items[i].Done = true
			items[i].UpdatedAt = time.Now().UTC()
			return items, nil
		}
	}
	return items, errors.New("checklist item not found")
}

// RemoveChecklistItem removes the checklist item with the given ID.
func RemoveChecklistItem(items []ChecklistItem, id string) ([]ChecklistItem, error) {
	for i, item := range items {
		if item.ID == id {
			return append(items[:i], items[i+1:]...), nil
		}
	}
	return items, errors.New("checklist item not found")
}

// PendingItems returns all incomplete checklist items for a snippet.
func PendingItems(items []ChecklistItem, snippetID string) []ChecklistItem {
	var result []ChecklistItem
	for _, item := range ChecklistFor(items, snippetID) {
		if !item.Done {
			result = append(result, item)
		}
	}
	return result
}
