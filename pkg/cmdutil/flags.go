package cmdutil

import (
	"sort"
	"strings"

	"github.com/pkg/errors"
)

// XorFlag will fail if 2 or more entries are true
func XorFlag(exclusiveFlags map[string]bool) error {
	fs := make([]string, 0)
	for n := range exclusiveFlags {
		if exclusiveFlags[n] {
			fs = append(fs, n)
		}
	}

	l := len(fs)
	if l < 2 {
		return nil
	}

	sort.Strings(fs)
	return FlagErrorWrap(errors.New(
		"the following flags can't be used together: " +
			strings.Join(fs[:l-1], ", ") + " and " + fs[l-1],
	))
}
