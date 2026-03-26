package util

import (
	"github.com/lucassabreu/clockify-cli/api"
)

// FillMissingBillableFn returns a step that derives the billable flag when
// the user did not explicitly set it, checking the task first, then the project.
func FillMissingBillableFn(c api.Client) Step {
	return func(dto TimeEntryDTO) (TimeEntryDTO, error) {
		if dto.Billable != nil || dto.ProjectID == "" {
			return dto, nil
		}

		if dto.TaskID != "" {
			t, err := c.GetTask(api.GetTaskParam{
				Workspace: dto.Workspace,
				ProjectID: dto.ProjectID,
				TaskID:    dto.TaskID,
			})
			if err != nil {
				return dto, err
			}
			b := t.Billable
			dto.Billable = &b
			return dto, nil
		}

		p, err := c.GetProject(api.GetProjectParam{
			Workspace: dto.Workspace,
			ProjectID: dto.ProjectID,
		})
		if err != nil {
			return dto, err
		}
		b := p.Billable
		dto.Billable = &b
		return dto, nil
	}
}

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
