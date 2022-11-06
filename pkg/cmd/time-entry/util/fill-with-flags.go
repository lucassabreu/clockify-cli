package util

import (
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/pflag"
)

// FillTimeEntryWithFlags will read the flags and fill the time entry with they
func FillTimeEntryWithFlags(flags *pflag.FlagSet) DoFn {
	return func(tei dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
		if flags.Changed("project") {
			p, _ := flags.GetString("project")
			if p != tei.ProjectID {
				tei.TaskID = ""
			}
			tei.ProjectID = p
		}

		if flags.Changed("description") {
			tei.Description, _ = flags.GetString("description")
		}

		if flags.Changed("task") {
			tei.TaskID, _ = flags.GetString("task")
		}

		if flags.Changed("tag") {
			tei.TagIDs, _ = flags.GetStringSlice("tag")
		}

		if flags.Changed("tags") {
			tei.TagIDs, _ = flags.GetStringSlice("tags")
		}

		if flags.Changed("not-billable") {
			b, _ := flags.GetBool("not-billable")
			tei.Billable = !b
		}

		var err error
		whenFlag := flags.Lookup("when")
		if whenFlag != nil && (whenFlag.Changed || whenFlag.DefValue != "") {
			whenString, _ := flags.GetString("when")
			var v time.Time
			if v, err = timehlp.ConvertToTime(whenString); err != nil {
				return tei, err
			}
			tei.TimeInterval.Start = v
		}

		if flags.Changed("end-at") {
			whenString, _ := flags.GetString("end-at")
			var v time.Time
			if v, err = timehlp.ConvertToTime(whenString); err != nil {
				return tei, err
			}
			tei.TimeInterval.End = &v
		}

		if flags.Changed("when-to-close") {
			whenString, _ := flags.GetString("when-to-close")
			var v time.Time
			if v, err = timehlp.ConvertToTime(whenString); err != nil {
				return tei, err
			}
			tei.TimeInterval.End = &v
		}

		return tei, nil
	}
}
