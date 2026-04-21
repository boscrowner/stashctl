package snippet

import (
	"time"
)

// Reminder represents a scheduled review reminder for a snippet.
type Reminder struct {
	SnippetID string    `json:"snippet_id"`
	DueAt     time.Time `json:"due_at"`
	Note      string    `json:"note,omitempty"`
}

// NewReminder creates a Reminder for the given snippet due at the specified time.
// Returns an error if snippetID is empty or dueAt is in the past.
func NewReminder(snippetID string, dueAt time.Time, note string) (Reminder, error) {
	if snippetID == "" {
		return Reminder{}, ErrEmptySnippetID
	}
	if !dueAt.After(time.Now()) {
		return Reminder{}, ErrDueAtInPast
	}
	return Reminder{
		SnippetID: snippetID,
		DueAt:     dueAt,
		Note:      note,
	}, nil
}

// DueReminders returns all reminders whose DueAt is on or before now.
func DueReminders(reminders []Reminder, now time.Time) []Reminder {
	var due []Reminder
	for _, r := range reminders {
		if !r.DueAt.After(now) {
			due = append(due, r)
		}
	}
	return due
}

// RemoveReminder removes the first reminder matching the given snippetID.
// Returns the updated slice and true if a reminder was removed.
func RemoveReminder(reminders []Reminder, snippetID string) ([]Reminder, bool) {
	for i, r := range reminders {
		if r.SnippetID == snippetID {
			return append(reminders[:i], reminders[i+1:]...), true
		}
	}
	return reminders, false
}

// Sentinel errors for reminder validation.
var (
	ErrEmptySnippetID = reminderError("snippet ID must not be empty")
	ErrDueAtInPast    = reminderError("due time must be in the future")
)

type reminderError string

func (e reminderError) Error() string { return string(e) }
