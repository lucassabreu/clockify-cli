package search

import (
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/pkg/errors"
)

func GetProjectByName(
	c api.Client,
	workspace,
	project string,
) (string, error) {
	id, err := findByName(project, "project", func() ([]named, error) {
		ps, err := c.GetProjects(api.GetProjectsParam{
			Workspace:       workspace,
			PaginationParam: api.AllPages(),
		})
		if err != nil {
			return []named{}, err
		}

		ns := make([]named, len(ps))
		for i := 0; i < len(ns); i++ {
			ns[i] = ps[i]
		}

		return ns, nil
	})

	if errors.Is(err, ErrEmptyReference) {
		return id, errors.New(
			"no project with id or name containing \"" +
				project + "\" was not found")
	}

	return id, err
}
