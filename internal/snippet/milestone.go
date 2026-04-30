package snippet

import (
	"errors"
	"time"
)

// Milestone represents a named goal attached to a snippet with an optional due date.
type Milestone struct {
	ID        string    `json:"id"`
	SnippetID string    `json:"snippet_id"`
	Name      string    `json:"name"`
	Note      string    `json:"note,omitempty"`
	Due       time.Time `json:"due,omitempty"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

const maxMilestoneNameLen = 80

// NewMilestone creates a new Milestone for the given snippet.
func NewMilestone(snippetID, name, note string, due time.Time) (Milestone, error) {
	if snippetID == "" {
		return Milestone{}, errors.New("milestone: snippet ID must not be empty")
	}
	if name == "" {
		return Milestone{}, errors.New("milestone: name must not be empty")
	}
	if len(name) > maxMilestoneNameLen {
		return Milestone{}, errors.New("milestone: name too long")
	}
	if !due.IsZero() && due.Before(time.Now()) {
		return Milestone{}, errors.New("milestone: due date must be in the future")
	}
	return Milestone{
		ID:        generateID(),
		SnippetID: snippetID,
		Name:      name,
		Note:      note,
		Due:       due,
		Done:      false,
		CreatedAt: time.Now(),
	}, nil
}

// MilestonesFor returns all milestones belonging to the given snippet.
func MilestonesFor(milestones []Milestone, snippetID string) []Milestone {
	var out []Milestone
	for _, m := range milestones {
		if m.SnippetID == snippetID {
			out = append(out, m)
		}
	}
	return out
}

// RemoveMilestone removes the milestone with the given ID from the slice.
func RemoveMilestone(milestones []Milestone, id string) ([]Milestone, error) {
	for i, m := range milestones {
		if m.ID == id {
			return append(milestones[:i], milestones[i+1:]...), nil
		}
	}
	return milestones, errors.New("milestone: not found")
}

// CompleteMilestone marks the milestone with the given ID as done.
func CompleteMilestone(milestones []Milestone, id string) ([]Milestone, error) {
	for i, m := range milestones {
		if m.ID == id {
			milestones[i].Done = true
			return milestones, nil
		}
	}
	return milestones, errors.New("milestone: not found")
}
