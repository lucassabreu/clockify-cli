package cmdutil

import (
	"sort"

	"github.com/lucassabreu/clockify-cli/strhlp"
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

	if len(fs) < 2 {
		return nil
	}

	sort.Strings(fs)
	fs = strhlp.Map(func(s string) string { return "`" + s + "`" }, fs)
	return FlagErrorWrap(errors.New(
		"the following flags can't be used together: " +
			strhlp.ListForHumans(fs),
	))
}
