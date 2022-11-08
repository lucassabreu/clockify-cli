package util

import (
	"time"

	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
)

type flagSet interface {
	Changed(string) bool
	GetString(string) (string, error)
	GetStringSlice(string) ([]string, error)
}

// FillTimeEntryWithFlags will read the flags and fill the time entry with they
func FillTimeEntryWithFlags(flags flagSet) Step {
	return func(dto TimeEntryDTO) (TimeEntryDTO, error) {
		if flags.Changed("project") {
			p, _ := flags.GetString("project")
			if p != dto.ProjectID {
				dto.TaskID = ""
			}
			dto.ProjectID = p
		}

		if flags.Changed("description") {
			dto.Description, _ = flags.GetString("description")
		}

		if flags.Changed("task") {
			dto.TaskID, _ = flags.GetString("task")
		}

		if flags.Changed("tag") {
			dto.TagIDs, _ = flags.GetStringSlice("tag")
		}

		if flags.Changed("tags") {
			dto.TagIDs, _ = flags.GetStringSlice("tags")
		}

		if flags.Changed("billable") {
			b := true
			dto.Billable = &b
		}

		if flags.Changed("not-billable") {
			b := false
			dto.Billable = &b
		}

		var err error
		if flags.Changed("when") {
			whenString, _ := flags.GetString("when")
			var v time.Time
			if v, err = timehlp.ConvertToTime(whenString); err != nil {
				return dto, err
			}
			dto.Start = v
		}

		if flags.Changed("when-to-close") {
			whenString, _ := flags.GetString("when-to-close")
			var v time.Time
			if v, err = timehlp.ConvertToTime(whenString); err != nil {
				return dto, err
			}
			dto.End = &v
		}

		return dto, nil
	}
}
