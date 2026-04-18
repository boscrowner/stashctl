package snippet

import (
	"sort"
	"strings"
)

// NormalizeTags lowercases and deduplicates a slice of tags,
// returning them in sorted order.
func NormalizeTags(tags []string) []string {
	seen := make(map[string]struct{}, len(tags))
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		t = strings.ToLower(strings.TrimSpace(t))
		if t == "" {
			continue
		}
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

// ParseTags splits a comma-separated tag string into a normalized slice.
func ParseTags(raw string) []string {
	parts := strings.Split(raw, ",")
	return NormalizeTags(parts)
}

// TagsEqual reports whether two tag slices contain the same tags
// regardless of order.
func TagsEqual(a, b []string) bool {
	na := NormalizeTags(a)
	nb := NormalizeTags(b)
	if len(na) != len(nb) {
		return false
	}
	for i := range na {
		if na[i] != nb[i] {
			return false
		}
	}
	return true
}
