package snippet_test

import (
	"strings"
	"testing"

	"github.com/user/stashctl/internal/snippet"
)

func makeReview(snippetID, reviewer string, status snippet.ReviewStatus) snippet.Review {
	r, err := snippet.NewReview(snippetID, reviewer, "looks good", status)
	if err != nil {
		panic(err)
	}
	return r
}

func TestNewReviewValid(t *testing.T) {
	r, err := snippet.NewReview("s1", "alice", "nice", snippet.ReviewApproved)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.ID == "" {
		t.Error("expected non-empty ID")
	}
	if r.SnippetID != "s1" {
		t.Errorf("expected snippet_id s1, got %s", r.SnippetID)
	}
}

func TestNewReviewEmptySnippetID(t *testing.T) {
	_, err := snippet.NewReview("", "alice", "note", snippet.ReviewPending)
	if err == nil {
		t.Error("expected error for empty snippet id")
	}
}

func TestNewReviewEmptyReviewer(t *testing.T) {
	_, err := snippet.NewReview("s1", "", "note", snippet.ReviewPending)
	if err == nil {
		t.Error("expected error for empty reviewer")
	}
}

func TestNewReviewCommentTooLong(t *testing.T) {
	long := strings.Repeat("x", 501)
	_, err := snippet.NewReview("s1", "bob", long, snippet.ReviewPending)
	if err == nil {
		t.Error("expected error for comment too long")
	}
}

func TestNewReviewUnknownStatus(t *testing.T) {
	_, err := snippet.NewReview("s1", "bob", "ok", snippet.ReviewStatus("unknown"))
	if err == nil {
		t.Error("expected error for unknown status")
	}
}

func TestReviewsFor(t *testing.T) {
	reviews := []snippet.Review{
		makeReview("s1", "alice", snippet.ReviewApproved),
		makeReview("s2", "bob", snippet.ReviewPending),
		makeReview("s1", "carol", snippet.ReviewRejected),
	}
	got := snippet.ReviewsFor(reviews, "s1")
	if len(got) != 2 {
		t.Errorf("expected 2 reviews, got %d", len(got))
	}
}

func TestRemoveReview(t *testing.T) {
	r := makeReview("s1", "alice", snippet.ReviewApproved)
	updated, err := snippet.RemoveReview([]snippet.Review{r}, r.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updated) != 0 {
		t.Errorf("expected empty slice, got %d", len(updated))
	}
}

func TestRemoveReviewNotFound(t *testing.T) {
	r := makeReview("s1", "alice", snippet.ReviewApproved)
	_, err := snippet.RemoveReview([]snippet.Review{r}, "nonexistent")
	if err == nil {
		t.Error("expected error for missing review")
	}
}

func TestApprovedReviews(t *testing.T) {
	reviews := []snippet.Review{
		makeReview("s1", "alice", snippet.ReviewApproved),
		makeReview("s1", "bob", snippet.ReviewPending),
		makeReview("s1", "carol", snippet.ReviewApproved),
	}
	got := snippet.ApprovedReviews(reviews)
	if len(got) != 2 {
		t.Errorf("expected 2 approved, got %d", len(got))
	}
}
