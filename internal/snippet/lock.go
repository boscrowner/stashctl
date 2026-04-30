package snippet

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Lock represents an exclusive edit lock on a snippet.
type Lock struct {
	ID        string    `json:"id"`
	SnippetID string    `json:"snippet_id"`
	Owner     string    `json:"owner"`
	AcquiredAt time.Time `json:"acquired_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// ErrLockConflict is returned when a snippet is already locked by another owner.
var ErrLockConflict = errors.New("snippet is already locked by another owner")

// ErrLockNotFound is returned when no lock exists for the given snippet.
var ErrLockNotFound = errors.New("lock not found")

// NewLock creates a new Lock for the given snippet and owner with the specified TTL.
func NewLock(snippetID, owner string, ttl time.Duration) (Lock, error) {
	if snippetID == "" {
		return Lock{}, errors.New("snippet ID must not be empty")
	}
	if owner == "" {
		return Lock{}, errors.New("owner must not be empty")
	}
	if ttl <= 0 {
		return Lock{}, errors.New("TTL must be positive")
	}
	now := time.Now().UTC()
	return Lock{
		ID:         uuid.NewString(),
		SnippetID:  snippetID,
		Owner:      owner,
		AcquiredAt: now,
		ExpiresAt:  now.Add(ttl),
	}, nil
}

// IsExpired reports whether the lock has passed its expiry time.
func (l Lock) IsExpired() bool {
	return time.Now().UTC().After(l.ExpiresAt)
}

// LocksFor returns all non-expired locks associated with a given snippet ID.
func LocksFor(locks []Lock, snippetID string) []Lock {
	var out []Lock
	for _, l := range locks {
		if l.SnippetID == snippetID && !l.IsExpired() {
			out = append(out, l)
		}
	}
	return out
}

// FindLock returns the first active lock for the snippet, if any.
func FindLock(locks []Lock, snippetID string) (Lock, bool) {
	for _, l := range locks {
		if l.SnippetID == snippetID && !l.IsExpired() {
			return l, true
		}
	}
	return Lock{}, false
}

// RemoveLock removes the lock with the given ID from the slice.
func RemoveLock(locks []Lock, id string) []Lock {
	out := locks[:0:0]
	for _, l := range locks {
		if l.ID != id {
			out = append(out, l)
		}
	}
	return out
}
