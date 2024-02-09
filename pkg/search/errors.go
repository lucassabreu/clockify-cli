package search

import (
	"sort"

	"github.com/lucassabreu/clockify-cli/strhlp"
)

// ErrNotFound represents a fail to identify a entity by its name or id
type ErrNotFound struct {
	EntityName string
	Reference  string
	Filters    map[string]string
}

func (e ErrNotFound) Error() string {
	sufix := ""
	if len(e.Filters) > 0 {
		sufix = " for "
		keys := make([]string, len(e.Filters))
		i := 0
		for k := range e.Filters {
			keys[i] = k
			i++
		}

		sort.Strings(keys)
		for i := range keys {
			keys[i] = keys[i] + " '" + e.Filters[keys[i]] + "'"
		}

		sufix = sufix + strhlp.ListForHumans(keys)
	}

	return "No " + e.EntityName + " with id or name containing '" +
		e.Reference + "' was found" + sufix
}
