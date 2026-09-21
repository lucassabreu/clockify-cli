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

// Candidate is an entity that matched an ambiguous reference
type Candidate struct {
	ID   string
	Name string
}

// maxCandidatesInError keeps the ambiguous error message readable when the
// reference matches a long list of entities
const maxCandidatesInError = 10

// ErrAmbiguous represents a reference that matched more than one entity, the
// caller must be more specific to identify which one it meant
type ErrAmbiguous struct {
	EntityName string
	Reference  string
	Candidates []Candidate
}

func (e ErrAmbiguous) Error() string {
	cs := make([]string, len(e.Candidates))
	for i, c := range e.Candidates {
		cs[i] = "'" + c.Name + "' (" + c.ID + ")"
	}

	return "More than one " + e.EntityName +
		" with id or name containing '" + e.Reference + "' was found: " +
		strhlp.LimitedListForHumans(cs, maxCandidatesInError)
}

func newErrAmbiguous(entityName, reference string, matches []named) ErrAmbiguous {
	cs := make([]Candidate, len(matches))
	for i, m := range matches {
		cs[i] = Candidate{ID: m.GetID(), Name: m.GetName()}
	}

	// the api listing order is not stable, sort so the message reads the
	// same across runs
	sort.Slice(cs, func(i, j int) bool {
		if cs[i].Name == cs[j].Name {
			return cs[i].ID < cs[j].ID
		}

		return cs[i].Name < cs[j].Name
	})

	return ErrAmbiguous{
		EntityName: entityName,
		Reference:  reference,
		Candidates: cs,
	}
}
