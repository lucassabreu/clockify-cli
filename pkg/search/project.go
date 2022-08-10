package search

import (
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
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

// GetProjectsByName will try to find projects containing the string on its
// name or id that matches the value
func GetProjectsByName(
	c api.Client,
	workspace string,
	projects []string,
) ([]string, error) {
	if len(projects) == 0 {
		return projects, nil
	}

	ts, err := c.GetProjects(api.GetProjectsParam{
		Workspace:       workspace,
		PaginationParam: api.AllPages(),
	})
	if err != nil {
		return projects, err
	}

	ns := make([]named, len(ts))
	for i := 0; i < len(ns); i++ {
		ns[i] = ts[i]
	}

	var g errgroup.Group
	for i := 0; i < len(projects); i++ {
		j := i
		g.Go(func() error {
			id, err := findByName(
				projects[j],
				"project", func() ([]named, error) { return ns, nil },
			)
			if err != nil {
				return err
			}

			projects[j] = id
			return nil
		})
	}

	return projects, g.Wait()
}
