package snippet

import (
	"sort"
	"strings"
)

// NormalizeTags lowercases, trims, deduplicates, and sorts tags.
func NormalizeTags(tags []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, t := range tags {
		t = strings.ToLower(strings.TrimSpace(t))
		if t == "" {
			continue
		}
		if _, ok := seen[t]; !ok {
			seen[t] = struct{}{}
			result = append(result, t)
		}
	}
	sort.Strings(result)
	return result
}

// ParseTags splits a comma-separated tag string into normalized tags.
func ParseTags(raw string) []string {
	parts := strings.Split(raw, ",")
	return NormalizeTags(parts)
}

// TagsEqual returns true if two tag slices contain the same tags.
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
