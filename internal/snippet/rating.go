package snippet

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Rating represents a user-assigned quality score for a snippet.
type Rating struct {
	ID        string    `json:"id"`
	SnippetID string    `json:"snippet_id"`
	Score     int       `json:"score"` // 1–5
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// NewRating creates a validated Rating for the given snippet.
func NewRating(snippetID string, score int, note string) (Rating, error) {
	if snippetID == "" {
		return Rating{}, errors.New("snippet_id is required")
	}
	if score < 1 || score > 5 {
		return Rating{}, errors.New("score must be between 1 and 5")
	}
	if len(note) > 280 {
		return Rating{}, errors.New("note must not exceed 280 characters")
	}
	return Rating{
		ID:        uuid.NewString(),
		SnippetID: snippetID,
		Score:     score,
		Note:      note,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// RatingsFor returns all ratings associated with the given snippet ID.
func RatingsFor(snippetID string, all []Rating) []Rating {
	var out []Rating
	for _, r := range all {
		if r.SnippetID == snippetID {
			out = append(out, r)
		}
	}
	return out
}

// AverageScore computes the mean score for a slice of ratings.
// Returns 0 if the slice is empty.
func AverageScore(ratings []Rating) float64 {
	if len(ratings) == 0 {
		return 0
	}
	sum := 0
	for _, r := range ratings {
		sum += r.Score
	}
	return float64(sum) / float64(len(ratings))
}

// RemoveRating removes the rating with the given ID from the slice.
func RemoveRating(id string, all []Rating) ([]Rating, error) {
	for i, r := range all {
		if r.ID == id {
			return append(all[:i], all[i+1:]...), nil
		}
	}
	return nil, errors.New("rating not found")
}
