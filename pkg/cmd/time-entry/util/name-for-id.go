package util

import (
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/search"
)

// GetAllowNameForIDsFn will try to find project/task/tags by their names if
// the value provided was not a ID
func GetAllowNameForIDsFn(config cmdutil.Config, c *api.Client) DoFn {
	if !config.GetBool(cmdutil.CONF_ALLOW_NAME_FOR_ID) {
		return nullCallback
	}

	cbs := []DoFn{
		lookupProject(c),
		lookupTask(c),
		lookupTags(c),
	}

	if config.IsInteractive() {
		cbs = disableErrorReporting(cbs)
	}

	return compose(cbs...)
}

func lookupProject(c *api.Client) DoFn {
	return func(te dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
		if te.ProjectID == "" {
			return te, nil
		}

		var err error
		te.ProjectID, err = search.GetProjectByName(
			c, te.WorkspaceID, te.ProjectID)
		return te, err
	}

}

func lookupTask(c *api.Client) DoFn {
	return func(te dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
		if te.TaskID == "" {
			return te, nil
		}

		var err error
		te.TaskID, err = search.GetTaskByName(
			c,
			api.GetTasksParam{
				Workspace: te.WorkspaceID,
				ProjectID: te.ProjectID,
				Active:    true,
			},
			te.TaskID)
		return te, err
	}
}

func lookupTags(c *api.Client) DoFn {
	return func(te dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
		if len(te.TagIDs) == 0 {
			return te, nil
		}

		var err error
		te.TagIDs, err = search.GetTagsByName(c, te.WorkspaceID, te.TagIDs)
		return te, err
	}

}

func disableErrorReporting(cbs []DoFn) []DoFn {
	for i := range cbs {
		cb := cbs[i]
		cbs[i] = func(tei dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
			tei, _ = cb(tei)
			return tei, nil
		}
	}
	return cbs
}
