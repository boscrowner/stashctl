package snippet_test

import (
	"strings"
	"testing"

	"github.com/user/stashctl/internal/snippet"
)

func makeDep(src, tgt, note string) snippet.Dependency {
	d, _ := snippet.NewDependency(src, tgt, note)
	return d
}

func TestNewDependencyValid(t *testing.T) {
	d, err := snippet.NewDependency("a", "b", "needs b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.SourceID != "a" || d.TargetID != "b" {
		t.Errorf("unexpected IDs: %+v", d)
	}
	if d.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestNewDependencyEmptySource(t *testing.T) {
	_, err := snippet.NewDependency("", "b", "")
	if err == nil {
		t.Error("expected error for empty source")
	}
}

func TestNewDependencyEmptyTarget(t *testing.T) {
	_, err := snippet.NewDependency("a", "", "")
	if err == nil {
		t.Error("expected error for empty target")
	}
}

func TestNewDependencySelfLoop(t *testing.T) {
	_, err := snippet.NewDependency("a", "a", "")
	if err == nil {
		t.Error("expected error for self-loop")
	}
}

func TestNewDependencyNoteTooLong(t *testing.T) {
	_, err := snippet.NewDependency("a", "b", strings.Repeat("x", 201))
	if err == nil {
		t.Error("expected error for note too long")
	}
}

func TestDependenciesFor(t *testing.T) {
	deps := []snippet.Dependency{
		makeDep("a", "b", ""),
		makeDep("a", "c", ""),
		makeDep("b", "c", ""),
	}
	result := snippet.DependenciesFor("a", deps)
	if len(result) != 2 {
		t.Errorf("expected 2, got %d", len(result))
	}
}

func TestDependentsOf(t *testing.T) {
	deps := []snippet.Dependency{
		makeDep("a", "c", ""),
		makeDep("b", "c", ""),
		makeDep("a", "b", ""),
	}
	result := snippet.DependentsOf("c", deps)
	if len(result) != 2 {
		t.Errorf("expected 2, got %d", len(result))
	}
}

func TestRemoveDependency(t *testing.T) {
	d := makeDep("a", "b", "")
	deps := []snippet.Dependency{d}
	updated, err := snippet.RemoveDependency(d.ID, deps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updated) != 0 {
		t.Errorf("expected empty slice, got %d", len(updated))
	}
}

func TestRemoveDependencyNotFound(t *testing.T) {
	deps := []snippet.Dependency{makeDep("a", "b", "")}
	_, err := snippet.RemoveDependency("nonexistent", deps)
	if err == nil {
		t.Error("expected error for missing dependency")
	}
}
