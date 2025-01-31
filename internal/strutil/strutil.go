package strutil

import (
	"slices"
	"sort"
	"unicode"
)

// Capitalize capitalizes the first character of a string.
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToTitle(runes[0])

	return string(runes)
}

// Dedup removes duplicates from a string slice.
func Dedup(s []string) []string {
	sort.Strings(s)
	return slices.Compact(s)
}
