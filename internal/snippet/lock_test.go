package snippet_test

import (
	"testing"
	"time"

	"github.com/user/stashctl/internal/snippet"
)

func makeLock(t *testing.T, snippetID, owner string, ttl time.Duration) snippet.Lock {
	t.Helper()
	l, err := snippet.NewLock(snippetID, owner, ttl)
	if err != nil {
		t.Fatalf("NewLock: %v", err)
	}
	return l
}

func TestNewLockValid(t *testing.T) {
	l := makeLock(t, "snip-1", "alice", time.Minute)
	if l.SnippetID != "snip-1" {
		t.Errorf("expected snippet_id snip-1, got %s", l.SnippetID)
	}
	if l.Owner != "alice" {
		t.Errorf("expected owner alice, got %s", l.Owner)
	}
	if l.ID == "" {
		t.Error("expected non-empty ID")
	}
	if l.IsExpired() {
		t.Error("newly created lock should not be expired")
	}
}

func TestNewLockEmptySnippetID(t *testing.T) {
	_, err := snippet.NewLock("", "alice", time.Minute)
	if err == nil {
		t.Error("expected error for empty snippet ID")
	}
}

func TestNewLockEmptyOwner(t *testing.T) {
	_, err := snippet.NewLock("snip-1", "", time.Minute)
	if err == nil {
		t.Error("expected error for empty owner")
	}
}

func TestNewLockNegativeTTL(t *testing.T) {
	_, err := snippet.NewLock("snip-1", "alice", -time.Second)
	if err == nil {
		t.Error("expected error for non-positive TTL")
	}
}

func TestIsExpired(t *testing.T) {
	l, _ := snippet.NewLock("snip-1", "alice", time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	if !l.IsExpired() {
		t.Error("lock should be expired after TTL")
	}
}

func TestLocksFor(t *testing.T) {
	l1 := makeLock(t, "snip-1", "alice", time.Minute)
	l2 := makeLock(t, "snip-2", "bob", time.Minute)
	l3, _ := snippet.NewLock("snip-1", "carol", time.Millisecond)
	time.Sleep(5 * time.Millisecond)

	result := snippet.LocksFor([]snippet.Lock{l1, l2, l3}, "snip-1")
	if len(result) != 1 || result[0].Owner != "alice" {
		t.Errorf("expected 1 active lock for alice, got %v", result)
	}
}

func TestFindLock(t *testing.T) {
	l := makeLock(t, "snip-1", "alice", time.Minute)
	found, ok := snippet.FindLock([]snippet.Lock{l}, "snip-1")
	if !ok {
		t.Fatal("expected to find lock")
	}
	if found.Owner != "alice" {
		t.Errorf("expected alice, got %s", found.Owner)
	}

	_, ok = snippet.FindLock([]snippet.Lock{l}, "snip-99")
	if ok {
		t.Error("expected no lock for unknown snippet")
	}
}

func TestRemoveLock(t *testing.T) {
	l1 := makeLock(t, "snip-1", "alice", time.Minute)
	l2 := makeLock(t, "snip-2", "bob", time.Minute)
	result := snippet.RemoveLock([]snippet.Lock{l1, l2}, l1.ID)
	if len(result) != 1 || result[0].ID != l2.ID {
		t.Errorf("expected only l2 to remain, got %v", result)
	}
}
