package util

import (
	"github.com/lucassabreu/clockify-cli/api"
)

// CreateTimeEntryFn will create a time entry
func CreateTimeEntryFn(c api.Client) Step {
	return func(dto TimeEntryDTO) (TimeEntryDTO, error) {
		te, err := c.CreateTimeEntry(api.CreateTimeEntryParam{
			Workspace:   dto.Workspace,
			Billable:    dto.Billable,
			Start:       dto.Start,
			End:         dto.End,
			ProjectID:   dto.ProjectID,
			Description: dto.Description,
			TagIDs:      dto.TagIDs,
			TaskID:      dto.TaskID,
		})

		if err != nil {
			return dto, err
		}

		return TimeEntryImplToDTO(te), nil
	}
}
