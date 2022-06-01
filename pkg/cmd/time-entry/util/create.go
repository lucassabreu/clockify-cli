package util

import (
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
)

// CreateTimeEntryFn will create a time entry
func CreateTimeEntryFn(c *api.Client) DoFn {
	return func(te dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
		return c.CreateTimeEntry(api.CreateTimeEntryParam{
			Workspace:   te.WorkspaceID,
			Billable:    te.Billable,
			Start:       te.TimeInterval.Start,
			End:         te.TimeInterval.End,
			ProjectID:   te.ProjectID,
			Description: te.Description,
			TagIDs:      te.TagIDs,
			TaskID:      te.TaskID,
		})
	}
}
