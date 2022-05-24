package search

import (
	"github.com/lucassabreu/clockify-cli/api"
)

func GetTaskByName(
	c *api.Client,
	workspace,
	project,
	task string,
) (string, error) {
	return findByName(task, "task", func() ([]named, error) {
		ts, err := c.GetTasks(api.GetTasksParam{
			Workspace:       workspace,
			ProjectID:       project,
			PaginationParam: api.AllPages(),
		})
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
