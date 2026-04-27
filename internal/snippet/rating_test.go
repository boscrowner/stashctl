package snippet_test

import (
	"strings"
	"testing"

	"github.com/user/stashctl/internal/snippet"
)

func makeRating(t *testing.T, snippetID string, score int, note string) snippet.Rating {
	t.Helper()
	r, err := snippet.NewRating(snippetID, score, note)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return r
}

func TestNewRatingValid(t *testing.T) {
	r := makeRating(t, "snip-1", 4, "pretty good")
	if r.ID == "" {
		t.Error("expected non-empty ID")
	}
	if r.Score != 4 {
		t.Errorf("expected score 4, got %d", r.Score)
	}
}

func TestNewRatingEmptySnippetID(t *testing.T) {
	_, err := snippet.NewRating("", 3, "")
	if err == nil {
		t.Error("expected error for empty snippet_id")
	}
}

func TestNewRatingScoreOutOfRange(t *testing.T) {
	for _, score := range []int{0, 6, -1} {
		_, err := snippet.NewRating("snip-1", score, "")
		if err == nil {
			t.Errorf("expected error for score %d", score)
		}
	}
}

func TestNewRatingNoteTooLong(t *testing.T) {
	long := strings.Repeat("x", 281)
	_, err := snippet.NewRating("snip-1", 3, long)
	if err == nil {
		t.Error("expected error for note exceeding 280 chars")
	}
}

func TestRatingsFor(t *testing.T) {
	r1 := makeRating(t, "snip-1", 5, "")
	r2 := makeRating(t, "snip-2", 3, "")
	r3 := makeRating(t, "snip-1", 2, "")

	got := snippet.RatingsFor("snip-1", []snippet.Rating{r1, r2, r3})
	if len(got) != 2 {
		t.Fatalf("expected 2 ratings, got %d", len(got))
	}
}

func TestAverageScore(t *testing.T) {
	r1 := makeRating(t, "s", 4, "")
	r2 := makeRating(t, "s", 2, "")
	avg := snippet.AverageScore([]snippet.Rating{r1, r2})
	if avg != 3.0 {
		t.Errorf("expected 3.0, got %f", avg)
	}
}

func TestAverageScoreEmpty(t *testing.T) {
	if avg := snippet.AverageScore(nil); avg != 0 {
		t.Errorf("expected 0 for empty ratings, got %f", avg)
	}
}

func TestRemoveRating(t *testing.T) {
	r := makeRating(t, "snip-1", 5, "")
	remaining, err := snippet.RemoveRating(r.ID, []snippet.Rating{r})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(remaining) != 0 {
		t.Errorf("expected empty slice after removal")
	}
}

func TestRemoveRatingNotFound(t *testing.T) {
	_, err := snippet.RemoveRating("nonexistent", []snippet.Rating{})
	if err == nil {
		t.Error("expected error for missing rating")
	}
}
