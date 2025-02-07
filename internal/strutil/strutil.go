package strutil

import (
	"slices"
	"sort"
)

// Dedup removes duplicates from a string slice.
func Dedup(s []string) []string {
	sort.Strings(s)
	return slices.Compact(s)
}
