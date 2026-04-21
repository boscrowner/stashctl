package snippet_test

import (
	"testing"
	"time"

	"github.com/user/stashctl/internal/snippet"
)

func makeReminder(t *testing.T, id string, offset time.Duration) snippet.Reminder {
	t.Helper()
	r, err := snippet.NewReminder(id, time.Now().Add(offset), "")
	if err != nil {
		t.Fatalf("makeReminder: %v", err)
	}
	return r
}

func TestNewReminderValid(t *testing.T) {
	_, err := snippet.NewReminder("abc", time.Now().Add(time.Hour), "review this")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNewReminderEmptyID(t *testing.T) {
	_, err := snippet.NewReminder("", time.Now().Add(time.Hour), "")
	if err == nil {
		t.Fatal("expected error for empty snippet ID")
	}
}

func TestNewReminderPastDue(t *testing.T) {
	_, err := snippet.NewReminder("abc", time.Now().Add(-time.Minute), "")
	if err == nil {
		t.Fatal("expected error for past due time")
	}
}

func TestDueRemindersReturnsOverdue(t *testing.T) {
	now := time.Now()
	reminders := []snippet.Reminder{
		{SnippetID: "a", DueAt: now.Add(-time.Hour)},
		{SnippetID: "b", DueAt: now.Add(time.Hour)},
		{SnippetID: "c", DueAt: now.Add(-time.Minute)},
	}
	due := snippet.DueReminders(reminders, now)
	if len(due) != 2 {
		t.Fatalf("expected 2 due reminders, got %d", len(due))
	}
}

func TestDueRemindersEmpty(t *testing.T) {
	due := snippet.DueReminders(nil, time.Now())
	if len(due) != 0 {
		t.Fatalf("expected empty result, got %d", len(due))
	}
}

func TestRemoveReminderFound(t *testing.T) {
	r := makeReminder(t, "x1", time.Hour)
	updated, ok := snippet.RemoveReminder([]snippet.Reminder{r}, "x1")
	if !ok {
		t.Fatal("expected reminder to be removed")
	}
	if len(updated) != 0 {
		t.Fatalf("expected empty slice, got %d", len(updated))
	}
}

func TestRemoveReminderNotFound(t *testing.T) {
	r := makeReminder(t, "x1", time.Hour)
	updated, ok := snippet.RemoveReminder([]snippet.Reminder{r}, "missing")
	if ok {
		t.Fatal("expected no removal")
	}
	if len(updated) != 1 {
		t.Fatalf("expected slice unchanged, got %d", len(updated))
	}
}
