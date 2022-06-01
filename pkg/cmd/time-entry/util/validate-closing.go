package util

import (
	"fmt"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
)

// ValidateClosingTimeEntry checks if the current time entry will fail to be
// stopped
func ValidateClosingTimeEntry(f cmdutil.Factory) DoFn {
	return func(tei dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
		c, err := f.Client()
		if err != nil {
			return tei, err
		}

		te, err := c.GetTimeEntryInProgress(api.GetTimeEntryInProgressParam{
			Workspace: tei.WorkspaceID,
			UserID:    tei.UserID,
		})

		if te == nil || err != nil {
			return tei, err
		}

		if err = validateTimeEntry(*te, f); err != nil {
			return tei, fmt.Errorf(
				"running time entry can't be ended: %w", err)
		}

		return tei, nil
	}
}
