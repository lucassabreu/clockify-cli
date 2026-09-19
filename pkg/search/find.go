package search

import (
	"errors"
	"strings"

	"github.com/lucassabreu/clockify-cli/strhlp"
)

type named interface {
	GetID() string
	GetName() string
}

var ErrEmptyReference = errors.New("no reference informed")

func findByName(
	r, entityName string, fn func() ([]named, error)) (string, error) {
	name := strhlp.Normalize(strings.TrimSpace(r))
	if name == "" {
		return r, ErrEmptyReference
	}

	l, err := fn()
	if err != nil {
		return r, err
	}

	for _, e := range l {
		if strings.ToLower(e.GetID()) == name {
			return e.GetID(), nil
		}
	}

	var exact, similar []named
	isSimilar := strhlp.IsSimilar(name)
	for _, e := range l {
		switch {
		case strhlp.Normalize(e.GetName()) == name:
			exact = append(exact, e)
		case isSimilar(e.GetName()):
			similar = append(similar, e)
		}
	}

	matches := exact
	if len(matches) == 0 {
		matches = similar
	}

	// the list may repeat an entity, e.g. the client of several projects,
	// repeated ids are not ambiguity
	matches = uniqueByID(matches)

	switch len(matches) {
	case 0:
		return r, ErrNotFound{
			EntityName: entityName,
			Reference:  r,
		}
	case 1:
		return matches[0].GetID(), nil
	default:
		return r, newErrAmbiguous(entityName, r, matches)
	}
}

func uniqueByID(l []named) []named {
	seen := make(map[string]bool, len(l))
	u := make([]named, 0, len(l))
	for _, e := range l {
		if seen[e.GetID()] {
			continue
		}

		seen[e.GetID()] = true
		u = append(u, e)
	}

	return u
}
