package util

import (
	"errors"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
)

// OutInProgressFn will stop the in progress time entry, if it exists
func OutInProgressFn(c api.Client) Step {
	return func(tei TimeEntryDTO) (TimeEntryDTO, error) {
		return tei, out(c, tei.Workspace, tei.UserID, tei.Start)
	}
}

func out(c api.Client, w, u string, end time.Time) error {
	if err := c.Out(api.OutParam{
		Workspace: w,
		UserID:    u,
		End:       end,
	}); getErrorCode(err) != 404 {
		return err
	}

	return nil
}

func getErrorCode(err error) int {
	var e dto.Error
	if errors.As(err, &e) {
		return e.Code
	}

	return 0
}
