package util

import (
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/search"
)

// GetAllowNameForIDsFn will try to find project/task/tags by their names if
// the value provided was not a ID
func GetAllowNameForIDsFn(config cmdutil.Config, c api.Client) Step {
	if !config.GetBool(cmdutil.CONF_ALLOW_NAME_FOR_ID) {
		return skip
	}

	cbs := []Step{
		lookupProject(c, config),
		lookupTask(c),
		lookupTags(c),
	}

	if config.IsInteractive() {
		cbs = disableErrorReporting(cbs)
	}

	return compose(cbs...)
}

func lookupProject(c api.Client, cnf cmdutil.Config) Step {
	return func(te TimeEntryDTO) (TimeEntryDTO, error) {
		if te.ProjectID == "" {
			return te, nil
		}

		var err error
		te.ProjectID, err = search.GetProjectByName(
			c, cnf, te.Workspace, te.ProjectID, te.Client)
		return te, err
	}

}

func lookupTask(c api.Client) Step {
	return func(te TimeEntryDTO) (TimeEntryDTO, error) {
		if te.TaskID == "" {
			return te, nil
		}

		var err error
		te.TaskID, err = search.GetTaskByName(
			c,
			api.GetTasksParam{
				Workspace: te.Workspace,
				ProjectID: te.ProjectID,
				Active:    true,
			},
			te.TaskID)
		return te, err
	}
}

func lookupTags(c api.Client) Step {
	return func(te TimeEntryDTO) (TimeEntryDTO, error) {
		if len(te.TagIDs) == 0 {
			return te, nil
		}

		var err error
		te.TagIDs, err = search.GetTagsByName(c, te.Workspace, true, te.TagIDs)
		return te, err
	}

}

func disableErrorReporting(cbs []Step) []Step {
	for i := range cbs {
		cb := cbs[i]
		cbs[i] = func(tei TimeEntryDTO) (TimeEntryDTO, error) {
			tei, _ = cb(tei)
			return tei, nil
		}
	}
	return cbs
}
