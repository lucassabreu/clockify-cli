package strhlp

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Normalize a string removes all non-english characters and lower
// case the string to help compare it with other strings
func Normalize(s string) string {
	if r, _, err := transform.String(
		transform.Chain(
			norm.NFD,
			runes.Remove(runes.In(unicode.Mn)),
			norm.NFC,
		),
		strings.ToLower(s),
	); err == nil {
		return r
	}

	return s
}

// InSlice will return true if the needle is one of the values in list
func InSlice(needle string, list []string) bool {
	return Search(needle, list) != -1
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

// Unique will remove all duplicated strings from the slice
func Unique(ss []string) []string {
	r := make([]string, 0)
	for _, s := range ss {
		if Search(s, r) == -1 {
			r = append(r, s)
		}
	}

	return r
}

// ListForHumans returns a string listing the strings from the parameter
//
// Example: ListForHumans([]string{"one", "two", "three"}) will output:
// "one, two and three"
func ListForHumans(s []string) string {
	if len(s) == 1 {
		return s[0]
	}

	return strings.Join(s[:len(s)-1], ", ") + " and " + s[len(s)-1]
}

// PadSpace will add spaces to the end of a string until it reaches the size
// set at the second parameter
func PadSpace(s string, size int) string {
	for i := len(s); i < size; i++ {
		s = s + " "
	}
	return s
}
