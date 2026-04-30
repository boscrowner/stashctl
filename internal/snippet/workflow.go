package snippet

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// WorkflowState represents a named stage in a snippet's lifecycle.
type WorkflowState string

const (
	WorkflowStateDraft    WorkflowState = "draft"
	WorkflowStateReview   WorkflowState = "review"
	WorkflowStateApproved WorkflowState = "approved"
	WorkflowStateArchived WorkflowState = "archived"
)

var knownWorkflowStates = map[WorkflowState]bool{
	WorkflowStateDraft:    true,
	WorkflowStateReview:   true,
	WorkflowStateApproved: true,
	WorkflowStateArchived: true,
}

// WorkflowTransition records a state change for a snippet.
type WorkflowTransition struct {
	ID        string        `json:"id"`
	SnippetID string        `json:"snippet_id"`
	From      WorkflowState `json:"from"`
	To        WorkflowState `json:"to"`
	Actor     string        `json:"actor"`
	Note      string        `json:"note"`
	At        time.Time     `json:"at"`
}

// NewWorkflowTransition creates a validated state transition record.
func NewWorkflowTransition(snippetID string, from, to WorkflowState, actor, note string) (WorkflowTransition, error) {
	if snippetID == "" {
		return WorkflowTransition{}, errors.New("snippet_id is required")
	}
	if actor == "" {
		return WorkflowTransition{}, errors.New("actor is required")
	}
	if !knownWorkflowStates[from] {
		return WorkflowTransition{}, errors.New("unknown source workflow state")
	}
	if !knownWorkflowStates[to] {
		return WorkflowTransition{}, errors.New("unknown target workflow state")
	}
	if from == to {
		return WorkflowTransition{}, errors.New("from and to states must differ")
	}
	if len(note) > 256 {
		return WorkflowTransition{}, errors.New("note exceeds 256 characters")
	}
	return WorkflowTransition{
		ID:        uuid.NewString(),
		SnippetID: snippetID,
		From:      from,
		To:        to,
		Actor:     actor,
		Note:      note,
		At:        time.Now().UTC(),
	}, nil
}

// TransitionsFor returns all transitions for a given snippet, ordered oldest first.
func TransitionsFor(snippetID string, all []WorkflowTransition) []WorkflowTransition {
	var out []WorkflowTransition
	for _, t := range all {
		if t.SnippetID == snippetID {
			out = append(out, t)
		}
	}
	return out
}

// CurrentState returns the most recent workflow state for a snippet.
// Returns WorkflowStateDraft as the default when no transitions exist.
func CurrentState(snippetID string, all []WorkflowTransition) WorkflowState {
	var latest WorkflowTransition
	for _, t := range all {
		if t.SnippetID == snippetID && (latest.ID == "" || t.At.After(latest.At)) {
			latest = t
		}
	}
	if latest.ID == "" {
		return WorkflowStateDraft
	}
	return latest.To
}
