package timeentry

import (
	"encoding/csv"
	"io"
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// TimeEntriesCSVPrint will print each time entry using the format string
func TimeEntriesCSVPrint(timeEntries []dto.TimeEntry, out io.Writer) error {
	w := csv.NewWriter(out)

	if err := w.Write([]string{
		"id",
		"description",
		"project.id",
		"project.name",
		"task.id",
		"task.name",
		"start",
		"end",
		"duration",
		"user.id",
		"user.email",
		"user.name",
		"tags...",
	}); err != nil {
		return err
	}

	format := func(t *time.Time) string {
		if t == nil {
			return ""
		}
		return t.In(time.Local).Format(TimeFormatFull)
	}

	for i := 0; i < len(timeEntries); i++ {
		te := timeEntries[i]
		var p dto.Project
		if te.Project != nil {
			p = *te.Project
		}

		end := time.Now()
		if te.TimeInterval.End != nil {
			end = *te.TimeInterval.End
		}

		if te.User == nil {
			u := dto.User{}
			te.User = &u
		}

		if te.Task == nil {
			t := dto.Task{}
			te.Task = &t
		}

		arr := []string{
			te.ID,
			te.Description,
			p.ID,
			p.Name,
			te.Task.ID,
			te.Task.Name,
			format(&te.TimeInterval.Start),
			format(te.TimeInterval.End),
			durationToString(end.Sub(te.TimeInterval.Start)),
			te.User.ID,
			te.User.Email,
			te.User.Name,
		}

		if err := w.Write(append(
			arr, tagsToStringSlice(te.Tags)...)); err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}
