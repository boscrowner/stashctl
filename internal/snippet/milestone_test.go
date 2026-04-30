package snippet

import (
	"testing"
	"time"
)

func makeMilestone(t *testing.T, snippetID, name string, due time.Time) Milestone {
	t.Helper()
	m, err := NewMilestone(snippetID, name, "", due)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return m
}

func TestNewMilestoneValid(t *testing.T) {
	due := time.Now().Add(24 * time.Hour)
	m, err := NewMilestone("s1", "v1.0 release", "first stable release", due)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if m.ID == "" {
		t.Error("expected non-empty ID")
	}
	if m.Done {
		t.Error("new milestone should not be done")
	}
}

func TestNewMilestoneEmptySnippetID(t *testing.T) {
	_, err := NewMilestone("", "goal", "", time.Time{})
	if err == nil {
		t.Error("expected error for empty snippet ID")
	}
}

func TestNewMilestoneEmptyName(t *testing.T) {
	_, err := NewMilestone("s1", "", "", time.Time{})
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestNewMilestoneNameTooLong(t *testing.T) {
	long := string(make([]byte, maxMilestoneNameLen+1))
	_, err := NewMilestone("s1", long, "", time.Time{})
	if err == nil {
		t.Error("expected error for name too long")
	}
}

func TestNewMilestonePastDue(t *testing.T) {
	past := time.Now().Add(-time.Hour)
	_, err := NewMilestone("s1", "old goal", "", past)
	if err == nil {
		t.Error("expected error for past due date")
	}
}

func TestMilestonesFor(t *testing.T) {
	m1 := makeMilestone(t, "s1", "alpha", time.Time{})
	m2 := makeMilestone(t, "s2", "beta", time.Time{})
	m3 := makeMilestone(t, "s1", "gamma", time.Time{})
	all := []Milestone{m1, m2, m3}
	got := MilestonesFor(all, "s1")
	if len(got) != 2 {
		t.Fatalf("expected 2 milestones, got %d", len(got))
	}
}

func TestRemoveMilestone(t *testing.T) {
	m := makeMilestone(t, "s1", "release", time.Time{})
	all := []Milestone{m}
	updated, err := RemoveMilestone(all, m.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updated) != 0 {
		t.Error("expected empty slice after removal")
	}
}

func TestRemoveMilestoneNotFound(t *testing.T) {
	all := []Milestone{}
	_, err := RemoveMilestone(all, "missing")
	if err == nil {
		t.Error("expected error for missing milestone")
	}
}

func TestCompleteMilestone(t *testing.T) {
	m := makeMilestone(t, "s1", "ship it", time.Time{})
	all := []Milestone{m}
	updated, err := CompleteMilestone(all, m.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !updated[0].Done {
		t.Error("expected milestone to be marked done")
	}
}
