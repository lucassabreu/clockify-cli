package uiutil

import (
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

	s, list := toList[dto.Tag](p.TagIDs, p.Tags)
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

	return listToEntities("tag", ids, p.Tags)
}
