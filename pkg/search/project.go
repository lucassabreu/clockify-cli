package search

import (
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func GetProjectByName(
	c api.Client,
	workspace string,
	project string,
	client string,
) (string, error) {
	ps, err := c.GetProjects(api.GetProjectsParam{
		Workspace:       workspace,
		PaginationParam: api.AllPages(),
	})
	if err != nil {
		return "", err
	}

	if ps, err = filterClientProjects(ps, client); err != nil {
		return "", err
	}

	id, err := findByName(project, "project", func() ([]named, error) {
		ns := make([]named, len(ps))
		for i := 0; i < len(ps); i++ {
			ns[i] = ps[i]
		}

		return ns, nil
	})

	var eNotFound ErrNotFound
	if errors.As(err, &eNotFound) {
		if client == "" {
			return id, err
		}

		return id, ErrNotFound{
			EntityName: eNotFound.EntityName,
			Reference:  eNotFound.Reference,
			Filters: map[string]string{
				"client": client,
			},
		}
	}

	return id, err
}

type namedStruct struct {
	ID   string
	Name string
}

func (c namedStruct) GetID() string {
	return c.ID
}

func (c namedStruct) GetName() string {
	return c.Name
}

func filterClientProjects(
	ps []dto.Project,
	client string,
) ([]dto.Project, error) {
	if client == "" {
		return ps, nil
	}

	clients := make([]named, len(ps))
	for i := 0; i < len(ps); i++ {
		clients[i] = namedStruct{
			ID:   ps[i].ClientID,
			Name: ps[i].ClientName,
		}
	}

	id, err := findByName(client, "client",
		func() ([]named, error) { return clients, nil })

	if err != nil {
		return ps, err
	}

	fPs := make([]dto.Project, 0)
	for i := 0; i < len(ps); i++ {
		if ps[i].ClientID != id {
			continue
		}

		fPs = append(fPs, ps[i])
	}

	return fPs, nil
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
