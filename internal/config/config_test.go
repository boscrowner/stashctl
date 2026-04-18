package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/stashctl/internal/config"
)

func TestDefaultHasValues(t *testing.T) {
	cfg, err := config.Default()
	if err != nil {
		t.Fatalf("Default() error: %v", err)
	}
	if cfg.StoreDir == "" {
		t.Error("expected non-empty StoreDir")
	}
	if cfg.DefaultFmt == "" {
		t.Error("expected non-empty DefaultFmt")
	}
}

func TestLoadMissingFileReturnsDefault(t *testing.T) {
	cfg, err := config.Load(filepath.Join(t.TempDir(), "config.json"))
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.DefaultFmt != "table" {
		t.Errorf("expected 'table', got %q", cfg.DefaultFmt)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	want := &config.Config{
		StoreDir:   "/tmp/snippets",
		DefaultFmt: "detail",
	}
	if err := config.Save(path, want); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	got, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if got.StoreDir != want.StoreDir {
		t.Errorf("StoreDir: got %q, want %q", got.StoreDir, want.StoreDir)
	}
	if got.DefaultFmt != want.DefaultFmt {
		t.Errorf("DefaultFmt: got %q, want %q", got.DefaultFmt, want.DefaultFmt)
	}
}

func TestSaveCreatesValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "config.json")
	cfg := &config.Config{StoreDir: "/tmp", DefaultFmt: "table"}
	if err := config.Save(path, cfg); err != nil {
		t.Fatalf("Save() error: %v", err)
	}
	data, _ := os.ReadFile(path)
	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Errorf("saved file is not valid JSON: %v", err)
	}
}
