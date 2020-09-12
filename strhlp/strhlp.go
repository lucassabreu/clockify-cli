package strhlp

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var t = transform.Chain(
	norm.NFD,
	runes.Remove(runes.In(unicode.Mn)),
	norm.NFC,
)

// Normalize a string removes all non-english characters and lower
// case the string to help compare it with other strings
func Normalize(s string) string {
	r, _, err := transform.String(t, strings.ToLower(s))
	if err != nil {
		return s
	}
	return r
}
