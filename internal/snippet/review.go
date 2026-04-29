package snippet

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ReviewStatus represents the state of a code review on a snippet.
type ReviewStatus string

const (
	ReviewPending  ReviewStatus = "pending"
	ReviewApproved ReviewStatus = "approved"
	ReviewRejected ReviewStatus = "rejected"
)

// Review represents a peer review attached to a snippet.
type Review struct {
	ID        string       `json:"id"`
	SnippetID string       `json:"snippet_id"`
	Reviewer  string       `json:"reviewer"`
	Status    ReviewStatus `json:"status"`
	Comment   string       `json:"comment"`
	CreatedAt time.Time    `json:"created_at"`
}

// NewReview creates a new Review for the given snippet.
func NewReview(snippetID, reviewer, comment string, status ReviewStatus) (Review, error) {
	if snippetID == "" {
		return Review{}, errors.New("snippet id is required")
	}
	if reviewer == "" {
		return Review{}, errors.New("reviewer is required")
	}
	if len(comment) > 500 {
		return Review{}, errors.New("comment must be 500 characters or fewer")
	}
	if status != ReviewPending && status != ReviewApproved && status != ReviewRejected {
		return Review{}, errors.New("unknown review status")
	}
	return Review{
		ID:        uuid.New().String(),
		SnippetID: snippetID,
		Reviewer:  reviewer,
		Status:    status,
		Comment:   comment,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// ReviewsFor returns all reviews for a given snippet ID.
func ReviewsFor(reviews []Review, snippetID string) []Review {
	var out []Review
	for _, r := range reviews {
		if r.SnippetID == snippetID {
			out = append(out, r)
		}
	}
	return out
}

// RemoveReview removes the review with the given ID from the slice.
func RemoveReview(reviews []Review, id string) ([]Review, error) {
	for i, r := range reviews {
		if r.ID == id {
			return append(reviews[:i], reviews[i+1:]...), nil
		}
	}
	return reviews, errors.New("review not found")
}

// ApprovedReviews returns only reviews with status approved.
func ApprovedReviews(reviews []Review) []Review {
	var out []Review
	for _, r := range reviews {
		if r.Status == ReviewApproved {
			out = append(out, r)
		}
	}
	return out
}
