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

	isSimilar := strhlp.IsSimilar(name)
	for _, e := range l {
		if strings.ToLower(e.GetID()) == name {
			return e.GetID(), nil
		}

		if isSimilar(e.GetName()) {
			return e.GetID(), nil
		}
	}

	return r, ErrNotFound{
		EntityName: entityName,
		Reference:  r,
	}
}
