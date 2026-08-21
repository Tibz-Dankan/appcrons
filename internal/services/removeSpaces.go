package services

import (
	"strings"
	"unicode"
)

// Removes all whitespace characters from the supplied string.
func RemoveSpaces(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, input)
}
