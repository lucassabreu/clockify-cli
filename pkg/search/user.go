package search

import (
	"github.com/lucassabreu/clockify-cli/api"
	"golang.org/x/sync/errgroup"
)

// GetUsersByName receives a list of id or names of clients and returns their
// ids
func GetUsersByName(
	c api.Client,
	workspace string,
	users []string,
) ([]string, error) {
	if len(users) == 0 {
		return users, nil
	}

	us, err := c.WorkspaceUsers(api.WorkspaceUsersParam{
		Workspace:       workspace,
		PaginationParam: api.AllPages(),
	})
	if err != nil {
		return users, err
	}

	ns := make([]named, len(us))
	for i := 0; i < len(ns); i++ {
		ns[i] = us[i]
	}

	var g errgroup.Group
	for i := 0; i < len(users); i++ {
		j := i
		g.Go(func() error {
			id, err := findByName(
				users[j], "user",
				func() ([]named, error) { return ns, nil },
			)
			if err != nil {
				return err
			}

			users[j] = id
			return nil
		})
	}

	return users, g.Wait()
}
