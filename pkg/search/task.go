package search

import (
	"github.com/lucassabreu/clockify-cli/api"
	"golang.org/x/sync/errgroup"
)

// GetTasksByName will try to find the first task containing the string on its
// name or id that matches the value
func GetTaskByName(
	c *api.Client,
	f api.GetTasksParam,
	task string,
) (string, error) {
	return findByName(task, "task", func() ([]named, error) {
		f.PaginationParam = api.AllPages()
		ts, err := c.GetTasks(f)
		if err != nil {
			return []named{}, err
		}

		ns := make([]named, len(ts))
		for i := 0; i < len(ns); i++ {
			ns[i] = ts[i]
		}

		return ns, nil
	})
}

// GetTasksByName will try to find tasks containing the string on its name or
// id that matches the value
func GetTasksByName(
	c *api.Client,
	f api.GetTasksParam,
	tasks []string,
) ([]string, error) {
	if len(tasks) == 0 {
		return tasks, nil
	}

	f.PaginationParam = api.AllPages()
	ts, err := c.GetTasks(f)
	if err != nil {
		return tasks, err
	}

	ns := make([]named, len(ts))
	for i := 0; i < len(ns); i++ {
		ns[i] = ts[i]
	}

	var g errgroup.Group
	for i := 0; i < len(tasks); i++ {
		j := i
		g.Go(func() error {
			id, err := findByName(
				tasks[j],
				"task", func() ([]named, error) { return ns, nil },
			)
			if err != nil {
				return err
			}

			tasks[j] = id
			return nil
		})
	}

	return tasks, g.Wait()
}
