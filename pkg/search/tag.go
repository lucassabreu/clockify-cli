package search

import (
	"github.com/lucassabreu/clockify-cli/api"
	"golang.org/x/sync/errgroup"
)

// GetTagsByName receives a list of id or names of tags and returns their ids
func GetTagsByName(
	c api.Client,
	workspace string,
	onlyActive bool,
	tags []string,
) ([]string, error) {
	if len(tags) == 0 {
		return tags, nil
	}
	var b *bool
	if onlyActive {
		f := false
		b = &f
	}

	ts, err := c.GetTags(api.GetTagsParam{
		Workspace:       workspace,
		Archived:        b,
		PaginationParam: api.AllPages(),
	})
	if err != nil {
		return tags, err
	}

	ns := make([]named, len(ts))
	for i := 0; i < len(ns); i++ {
		ns[i] = ts[i]
	}

	var g errgroup.Group
	for i := 0; i < len(tags); i++ {
		j := i
		g.Go(func() error {
			id, err := findByName(
				tags[j],
				"tag", func() ([]named, error) { return ns, nil },
			)
			if err != nil {
				return err
			}

			tags[j] = id
			return nil
		})
	}

	return tags, g.Wait()
}
