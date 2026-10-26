package analyze

import (
	"strings"
	"unicode"
)

// Lowercases text and splits it into terms.
func Tokenize(text string) []string {
	lowercase := strings.ToLower(text)
	split := strings.FieldsFunc(lowercase, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	return split
}
