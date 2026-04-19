package main

import (
	"strings"
	"testing"
)

func TestStatsEmpty(t *testing.T) {
	env := setupTest(t)
	out, err := env.run("stats")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Total snippets: 0") {
		t.Fatalf("expected zero total, got: %s", out)
	}
}

func TestStatsWithSnippets(t *testing.T) {
	env := setupTest(t)
	_, err := env.run("add", "--title", "Hello Go", "--content", "fmt.Println()", "--language", "go", "--tags", "cli,util")
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}
	_, err = env.run("add", "--title", "Hello Python", "--content", "print()", "--language", "python", "--tags", "util")
	if err != nil {
		t.Fatalf("add failed: %v", err)
	}
	out, err := env.run("stats")
	if err != nil {
		t.Fatalf("stats failed: %v", err)
	}
	if !strings.Contains(out, "Total snippets: 2") {
		t.Fatalf("expected 2 total, got: %s", out)
	}
	if !strings.Contains(out, "go") {
		t.Fatalf("expected go in output: %s", out)
	}
	if !strings.Contains(out, "python") {
		t.Fatalf("expected python in output: %s", out)
	}
	if !strings.Contains(out, "util") {
		t.Fatalf("expected util tag in output: %s", out)
	}
}
