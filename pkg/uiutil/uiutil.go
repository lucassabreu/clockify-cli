package uiutil

import (
	"errors"
	"strings"
)

type named interface {
	GetName() string
	GetID() string
}

func toList[T named](
	d []string, options []T) (strD []string, strOpts []string) {
	strOpts = make([]string, len(options))
	for i := range options {
		strOpts[i] = options[i].GetID() + " - " + options[i].GetName()
	}

	strD = make([]string, len(d))
	for i := range d {
		for _, o := range strOpts {
			if strings.HasPrefix(o, d[i]) {
				strD[i] = o
			}
		}
	}
	return
}

func listToEntities[T named](
	name string, ids []string, entities []T) ([]T, error) {
	selected := make([]T, len(ids))
	for i, t := range ids {
		found := false
		t = strings.TrimSpace(t[0:strings.Index(t, " - ")])
		for j := range entities {
			if entities[j].GetID() == t {
				selected[i] = entities[j]
				found = true
			}
		}

		if !found {
			return []T{}, errors.New(
				name + ` with id "` + t + `" not found`)
		}
	}

	return selected, nil
}
