package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/user/stashctl/internal/snippet"
)

const workspaceFile = "workspaces.json"

// SaveWorkspace persists a workspace to disk.
func (s *Store) SaveWorkspace(w *snippet.Workspace) error {
	workspaces, err := s.loadWorkspaces()
	if err != nil {
		return err
	}
	for i, existing := range workspaces {
		if existing.ID == w.ID {
			workspaces[i] = w
			return s.writeWorkspaces(workspaces)
		}
	}
	workspaces = append(workspaces, w)
	return s.writeWorkspaces(workspaces)
}

// GetWorkspace retrieves a workspace by ID.
func (s *Store) GetWorkspace(id string) (*snippet.Workspace, error) {
	workspaces, err := s.loadWorkspaces()
	if err != nil {
		return nil, err
	}
	for _, w := range workspaces {
		if w.ID == id {
			return w, nil
		}
	}
	return nil, errors.New("workspace not found: " + id)
}

// ListWorkspaces returns all stored workspaces.
func (s *Store) ListWorkspaces() ([]*snippet.Workspace, error) {
	return s.loadWorkspaces()
}

// DeleteWorkspace removes a workspace by ID.
func (s *Store) DeleteWorkspace(id string) error {
	workspaces, err := s.loadWorkspaces()
	if err != nil {
		return err
	}
	for i, w := range workspaces {
		if w.ID == id {
			workspaces = append(workspaces[:i], workspaces[i+1:]...)
			return s.writeWorkspaces(workspaces)
		}
	}
	return errors.New("workspace not found: " + id)
}

func (s *Store) workspacePath() string {
	return filepath.Join(s.Dir(), workspaceFile)
}

func (s *Store) loadWorkspaces() ([]*snippet.Workspace, error) {
	data, err := os.ReadFile(s.workspacePath())
	if errors.Is(err, os.ErrNotExist) {
		return []*snippet.Workspace{}, nil
	}
	if err != nil {
		return nil, err
	}
	var workspaces []*snippet.Workspace
	if err := json.Unmarshal(data, &workspaces); err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (s *Store) writeWorkspaces(workspaces []*snippet.Workspace) error {
	data, err := json.MarshalIndent(workspaces, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.workspacePath(), data, 0o644)
}
