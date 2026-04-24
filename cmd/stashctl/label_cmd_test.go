package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestLabelCreateAndList(t *testing.T) {
	env := setupTest(t)

	// create
	out, err := env.run("label", "create", "urgent", "--color", "#ff0000")
	if err != nil {
		t.Fatalf("create: %v\nout: %s", err, out)
	}
	if !strings.Contains(out, "urgent") {
		t.Errorf("expected 'urgent' in output, got: %s", out)
	}

	// list
	out, err = env.run("label", "list")
	if err != nil {
		t.Fatalf("list: %v\nout: %s", err, out)
	}
	if !strings.Contains(out, "urgent") {
		t.Errorf("expected 'urgent' in list output, got: %s", out)
	}
	if !strings.Contains(out, "#ff0000") {
		t.Errorf("expected color '#ff0000' in list output, got: %s", out)
	}
}

func TestLabelListEmpty(t *testing.T) {
	env := setupTest(t)

	out, err := env.run("label", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "no labels") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestLabelRemove(t *testing.T) {
	env := setupTest(t)

	if _, err := env.run("label", "create", "temp"); err != nil {
		t.Fatalf("create: %v", err)
	}

	out, err := env.run("label", "remove", "temp")
	if err != nil {
		t.Fatalf("remove: %v\nout: %s", err, out)
	}
	if !strings.Contains(out, "removed") {
		t.Errorf("expected 'removed' in output, got: %s", out)
	}

	// list should be empty again
	out, _ = env.run("label", "list")
	if !strings.Contains(out, "no labels") {
		t.Errorf("expected empty list after removal, got: %s", out)
	}
}

func TestLabelRemoveNotFound(t *testing.T) {
	env := setupTest(t)

	_, err := env.run("label", "remove", "ghost")
	if err == nil {
		t.Error("expected error for missing label")
	}
}

func TestLabelCreateInvalidName(t *testing.T) {
	env := setupTest(t)

	var buf bytes.Buffer
	_ = buf // suppress unused warning
	_, err := env.run("label", "create", strings.Repeat("x", 33))
	if err == nil {
		t.Error("expected error for name too long")
	}
}
