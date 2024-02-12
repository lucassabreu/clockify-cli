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

	if p.Tags == nil || len(p.Tags) == 0 {
		return []dto.Tag{}, nil
	}

	if p.Message == "" {
		p.Message = "Choose your tags:"
	}

	list := make([]string, len(p.Tags))
	for i := range p.Tags {
		list[i] = p.Tags[i].ID + " - " + p.Tags[i].Name
	}

	for i := range p.TagIDs {
		for _, id := range list {
			if strings.HasPrefix(id, p.TagIDs[i]) {
				p.TagIDs[i] = id
			}
		}
	}

	v := func(s []string) error { return nil }
	if p.Force {
		v = func(s []string) error {
			if len(s) == 0 {
				return errors.New("at least one tag should be selected")
			}
			return nil
		}
	}

	ids, err := p.UI.AskManyFromOptions(p.Message, list, p.TagIDs, v)
	if err != nil || len(ids) == 0 {
		return []dto.Tag{}, err
	}

	tags := make([]dto.Tag, len(ids))
	for _, t := range ids {
		found := false
		t = strings.TrimSpace(t[0:strings.Index(t, " - ")])
		for i := range p.Tags {
			if p.Tags[i].ID == t {
				tags[i] = p.Tags[i]
				found = true
			}
		}

		if !found {
			return []dto.Tag{}, errors.New(`tag with id "` + t + `" not found`)
		}
	}

	return tags, nil
}
