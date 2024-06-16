package uiutil

import (
	"strings"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/ui"
	"github.com/pkg/errors"
)

// AskTagsParam informs what options to display while asking for a tag
type AskTagsParam struct {
	UI      ui.UI
	TagIDs  []string
	Tags    []dto.Tag
	Message string
	Force   bool
}

// AskTags asks the user for a tag from options
func AskTags(p AskTagsParam) ([]dto.Tag, error) {
	if p.UI == nil {
		return nil, errors.New("UI must be informed")
	}

	if p.Message == "" {
		p.Message = "Choose your tags:"
	}

	s, list := tagsToList(p.TagIDs, p.Tags)
	v := func(s []string) error { return nil }
	if p.Force {
		v = func(s []string) error {
			if len(s) == 0 {
				return errors.New("at least one tag should be selected")
			}
			return nil
		}
	}

	ids, err := p.UI.AskManyFromOptions(p.Message, list, s, v)
	if err != nil || len(ids) == 0 {
		return []dto.Tag{}, err
	}

	return listToTags("tag", ids, p.Tags)
}
func tagsToList(
	d []string, options []dto.Tag) (strD []string, strOpts []string) {
	strOpts = make([]string, len(options))
	for i := range options {
		strOpts[i] = options[i].ID + " - " + options[i].Name
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

func listToTags(
	name string, ids []string, entities []dto.Tag) ([]dto.Tag, error) {
	selected := make([]dto.Tag, len(ids))
	for i, t := range ids {
		found := false
		t = strings.TrimSpace(t[0:strings.Index(t, " - ")])
		for j := range entities {
			if entities[j].ID == t {
				selected[i] = entities[j]
				found = true
			}
		}

		if !found {
			return []dto.Tag{}, errors.New(
				name + ` with id "` + t + `" not found`)
		}
	}

	return selected, nil
}
