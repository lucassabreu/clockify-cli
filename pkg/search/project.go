package search

import (
	"github.com/lucassabreu/clockify-cli/api"
)

func GetProjectByName(
	c *api.Client,
	workspace,
	project string,
) (string, error) {
	return findByName(project, "project", func() ([]named, error) {
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
}
