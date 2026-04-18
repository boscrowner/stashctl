package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/stashctl/internal/config"
	"github.com/user/stashctl/internal/store"
)

func setupTest(t *testing.T) (*config.Config, *store.Store) {
	t.Helper()
	dir := t.TempDir()
	cfg := config.Default()
	cfg.StorePath = filepath.Join(dir, "snippets.json")
	cfg.ColorOutput = false
	st, err := store.New(cfg.StorePath)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return cfg, st
}

func TestAddAndListCmd(t *testing.T) {
	cfg, st := setupTest(t)
	root := newRootCmd(cfg, st)

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(os.Stderr)

	root.SetArgs([]string{"add", "hello world", "fmt.Println()", "--lang", "go", "--tags", "go,print"})
	if err := root.Execute(); err != nil {
		t.Fatalf("add cmd: %v", err)
	}

	buf.Reset()
	root.SetArgs([]string{"list"})
	if err := root.Execute(); err != nil {
		t.Fatalf("list cmd: %v", err)
	}
}

func TestSearchCmd(t *testing.T) {
	cfg, st := setupTest(t)
	root := newRootCmd(cfg, st)
	root.SetErr(os.Stderr)

	root.SetArgs([]string{"add", "hello world", "fmt.Println()", "--lang", "go"})
	if err := root.Execute(); err != nil {
		t.Fatalf("add: %v", err)
	}

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"search", "hello"})
	if err := root.Execute(); err != nil {
		t.Fatalf("search: %v", err)
	}
}

func TestExportCmd(t *testing.T) {
	cfg, st := setupTest(t)
	root := newRootCmd(cfg, st)
	root.SetErr(os.Stderr)

	root.SetArgs([]string{"add", "my snippet", "body", "--lang", "go"})
	_ = root.Execute()

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"export", "--format", "json"})
	if err := root.Execute(); err != nil {
		t.Fatalf("export: %v", err)
	}
	if !strings.Contains(buf.String(), "my snippet") {
		t.Errorf("expected snippet in export output, got: %s", buf.String())
	}
}
