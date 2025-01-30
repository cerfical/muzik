package strutil

import "unicode"

// Capitalize capitalizes the first character of a string.
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToTitle(runes[0])

	return string(runes)
}
