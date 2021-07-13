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

// Search will search for a exact match of the string on the slide
// provided and if found will return its index position, or -1 if not
func Search(s string, list []string) int {
	for i, sl := range list {
		if sl == s {
			return i
		}
	}

	return -1
}

// Map will apply the map function provided on every entry of
// the string slice and return a slice with the changes
func Map(f func(string) string, s []string) []string {
	for i, e := range s {
		s[i] = f(e)
	}
	return s
}

// Filter will run the function for every entry of the slice
// and will return a new slice with the entries that return true
// when used on the function
func Filter(f func(string) bool, s []string) []string {
	ns := make([]string, 0)
	for _, e := range s {
		if f(e) {
			ns = append(ns, e)
		}
	}
	return ns
}

// Merge will combine all string slices into a single one
func Merge(ss ...[]string) []string {
	s := make([]string, 0)
	for i := range ss {
		s = append(s, ss[i]...)
	}
	return s
}
