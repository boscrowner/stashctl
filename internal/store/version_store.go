package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/user/stashctl/internal/snippet"
)

const versionsFile = "versions.json"

func (s *Store) versionsPath() string {
	return filepath.Join(s.dir, versionsFile)
}

func (s *Store) LoadVersions() ([]snippet.Version, error) {
	data, err := os.ReadFile(s.versionsPath())
	if errors.Is(err, os.ErrNotExist) {
		return []snippet.Version{}, nil
	}
	if err != nil {
		return nil, err
	}
	var versions []snippet.Version
	if err := json.Unmarshal(data, &versions); err != nil {
		return nil, err
	}
	return versions, nil
}

func (s *Store) SaveVersions(versions []snippet.Version) error {
	data, err := json.MarshalIndent(versions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.versionsPath(), data, 0644)
}

func (s *Store) AddVersion(v snippet.Version) error {
	versions, err := s.LoadVersions()
	if err != nil {
		return err
	}
	versions = append(versions, v)
	return s.SaveVersions(versions)
}

func (s *Store) ListVersions(snippetID string) ([]snippet.Version, error) {
	versions, err := s.LoadVersions()
	if err != nil {
		return nil, err
	}
	return snippet.VersionsFor(snippetID, versions), nil
}

func (s *Store) DeleteVersionsFor(snippetID string) error {
	versions, err := s.LoadVersions()
	if err != nil {
		return err
	}
	var kept []snippet.Version
	for _, v := range versions {
		if v.SnippetID != snippetID {
			kept = append(kept, v)
		}
	}
	return s.SaveVersions(kept)
}
