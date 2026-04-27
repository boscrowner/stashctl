package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/user/stashctl/internal/snippet"
)

const dependencyFile = "dependencies.json"

func (s *Store) dependencyPath() string {
	return filepath.Join(s.dir, dependencyFile)
}

// LoadDependencies reads all dependencies from disk.
func (s *Store) LoadDependencies() ([]snippet.Dependency, error) {
	data, err := os.ReadFile(s.dependencyPath())
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var deps []snippet.Dependency
	if err := json.Unmarshal(data, &deps); err != nil {
		return nil, err
	}
	return deps, nil
}

// SaveDependencies persists all dependencies to disk.
func (s *Store) SaveDependencies(deps []snippet.Dependency) error {
	data, err := json.MarshalIndent(deps, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.dependencyPath(), data, 0o644)
}

// AddDependency appends a new dependency and saves.
func (s *Store) AddDependency(dep snippet.Dependency) error {
	deps, err := s.LoadDependencies()
	if err != nil {
		return err
	}
	deps = append(deps, dep)
	return s.SaveDependencies(deps)
}

// DeleteDependency removes a dependency by ID and saves.
func (s *Store) DeleteDependency(id string) error {
	deps, err := s.LoadDependencies()
	if err != nil {
		return err
	}
	updated, err := snippet.RemoveDependency(id, deps)
	if err != nil {
		return err
	}
	return s.SaveDependencies(updated)
}
