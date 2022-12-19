package util

import (
	"fmt"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
)

// ValidateClosingTimeEntry checks if the current time entry will fail to be
// stopped
func ValidateClosingTimeEntry(f cmdutil.Factory) Step {
	return func(dto TimeEntryDTO) (TimeEntryDTO, error) {
		c, err := f.Client()
		if err != nil {
			return dto, err
		}

		te, err := c.GetTimeEntryInProgress(api.GetTimeEntryInProgressParam{
			Workspace: dto.Workspace,
			UserID:    dto.UserID,
		})

		if te == nil || err != nil {
			return dto, err
		}

		if err = validateTimeEntry(TimeEntryImplToDTO(*te), f); err != nil {
			return dto, fmt.Errorf(
				"running time entry can't be ended: %w", err)
		}

		return dto, nil
	}
}
