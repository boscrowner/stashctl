package snippet_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/stashctl/internal/snippet"
)

func makeTransition(snippetID string, from, to snippet.WorkflowState, actor string, at time.Time) snippet.WorkflowTransition {
	t, err := snippet.NewWorkflowTransition(snippetID, from, to, actor, "")
	if err != nil {
		panic(err)
	}
	t.At = at
	return t
}

func TestNewWorkflowTransitionValid(t *testing.T) {
	tr, err := snippet.NewWorkflowTransition("s1", snippet.WorkflowStateDraft, snippet.WorkflowStateReview, "alice", "ready")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.ID == "" {
		t.Error("expected non-empty ID")
	}
	if tr.To != snippet.WorkflowStateReview {
		t.Errorf("expected To=review, got %s", tr.To)
	}
}

func TestNewWorkflowTransitionEmptySnippetID(t *testing.T) {
	_, err := snippet.NewWorkflowTransition("", snippet.WorkflowStateDraft, snippet.WorkflowStateReview, "alice", "")
	if err == nil {
		t.Fatal("expected error for empty snippet_id")
	}
}

func TestNewWorkflowTransitionEmptyActor(t *testing.T) {
	_, err := snippet.NewWorkflowTransition("s1", snippet.WorkflowStateDraft, snippet.WorkflowStateReview, "", "")
	if err == nil {
		t.Fatal("expected error for empty actor")
	}
}

func TestNewWorkflowTransitionSameState(t *testing.T) {
	_, err := snippet.NewWorkflowTransition("s1", snippet.WorkflowStateDraft, snippet.WorkflowStateDraft, "alice", "")
	if err == nil {
		t.Fatal("expected error when from == to")
	}
}

func TestNewWorkflowTransitionUnknownState(t *testing.T) {
	_, err := snippet.NewWorkflowTransition("s1", "unknown", snippet.WorkflowStateReview, "alice", "")
	if err == nil {
		t.Fatal("expected error for unknown source state")
	}
}

func TestNewWorkflowTransitionNoteTooLong(t *testing.T) {
	_, err := snippet.NewWorkflowTransition("s1", snippet.WorkflowStateDraft, snippet.WorkflowStateReview, "alice", strings.Repeat("x", 257))
	if err == nil {
		t.Fatal("expected error for note exceeding 256 chars")
	}
}

func TestTransitionsFor(t *testing.T) {
	now := time.Now()
	all := []snippet.WorkflowTransition{
		makeTransition("s1", snippet.WorkflowStateDraft, snippet.WorkflowStateReview, "alice", now),
		makeTransition("s2", snippet.WorkflowStateDraft, snippet.WorkflowStateApproved, "bob", now),
		makeTransition("s1", snippet.WorkflowStateReview, snippet.WorkflowStateApproved, "carol", now.Add(time.Minute)),
	}
	result := snippet.TransitionsFor("s1", all)
	if len(result) != 2 {
		t.Errorf("expected 2 transitions for s1, got %d", len(result))
	}
}

func TestCurrentStateDefault(t *testing.T) {
	state := snippet.CurrentState("s99", nil)
	if state != snippet.WorkflowStateDraft {
		t.Errorf("expected draft as default, got %s", state)
	}
}

func TestCurrentStateLatest(t *testing.T) {
	now := time.Now()
	all := []snippet.WorkflowTransition{
		makeTransition("s1", snippet.WorkflowStateDraft, snippet.WorkflowStateReview, "alice", now),
		makeTransition("s1", snippet.WorkflowStateReview, snippet.WorkflowStateApproved, "bob", now.Add(time.Hour)),
	}
	state := snippet.CurrentState("s1", all)
	if state != snippet.WorkflowStateApproved {
		t.Errorf("expected approved, got %s", state)
	}
}
