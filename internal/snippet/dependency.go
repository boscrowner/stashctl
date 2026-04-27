package snippet

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Dependency represents a directional link from one snippet to another,
// indicating that the source snippet depends on the target snippet.
type Dependency struct {
	ID          string    `json:"id"`
	SourceID    string    `json:"source_id"`
	TargetID    string    `json:"target_id"`
	Note        string    `json:"note,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// NewDependency creates a new Dependency linking sourceID -> targetID.
func NewDependency(sourceID, targetID, note string) (Dependency, error) {
	if sourceID == "" {
		return Dependency{}, errors.New("source_id is required")
	}
	if targetID == "" {
		return Dependency{}, errors.New("target_id is required")
	}
	if sourceID == targetID {
		return Dependency{}, errors.New("source and target must differ")
	}
	if len(note) > 200 {
		return Dependency{}, fmt.Errorf("note exceeds 200 characters")
	}
	return Dependency{
		ID:        uuid.NewString(),
		SourceID:  sourceID,
		TargetID:  targetID,
		Note:      note,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// DependenciesFor returns all dependencies where the given snippet is the source.
func DependenciesFor(snippetID string, all []Dependency) []Dependency {
	var out []Dependency
	for _, d := range all {
		if d.SourceID == snippetID {
			out = append(out, d)
		}
	}
	return out
}

// DependentsOf returns all dependencies where the given snippet is the target.
func DependentsOf(snippetID string, all []Dependency) []Dependency {
	var out []Dependency
	for _, d := range all {
		if d.TargetID == snippetID {
			out = append(out, d)
		}
	}
	return out
}

// RemoveDependency removes a dependency by ID.
func RemoveDependency(id string, all []Dependency) ([]Dependency, error) {
	for i, d := range all {
		if d.ID == id {
			return append(all[:i], all[i+1:]...), nil
		}
	}
	return all, fmt.Errorf("dependency %q not found", id)
}
