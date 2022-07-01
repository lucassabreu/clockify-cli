package search

import (
	"github.com/lucassabreu/clockify-cli/api"
	"golang.org/x/sync/errgroup"
)

// GetClientsByName receives a list of id or names of clients and returns their
// ids
func GetClientsByName(
	c *api.Client,
	workspace string,
	clients []string,
) ([]string, error) {
	if len(clients) == 0 {
		return clients, nil
	}

	cs, err := c.GetClients(api.GetClientsParam{
		Workspace:       workspace,
		PaginationParam: api.AllPages(),
	})
	if err != nil {
		return clients, err
	}

	ns := make([]named, len(cs))
	for i := 0; i < len(ns); i++ {
		ns[i] = cs[i]
	}

	var g errgroup.Group
	for i := 0; i < len(clients); i++ {
		j := i
		g.Go(func() error {
			id, err := findByName(
				clients[j],
				"client", func() ([]named, error) { return ns, nil },
			)
			if err != nil {
				return err
			}

			clients[j] = id
			return nil
		})
	}

	return clients, g.Wait()
}
