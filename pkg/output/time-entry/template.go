package timeentry

import (
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/output/util"
)

// TimeEntriesPrintWithTemplate will print each time entry using the format
// string
func TimeEntriesPrintWithTemplate(
	format string,
) func([]dto.TimeEntry, io.Writer) error {
	return func(timeEntries []dto.TimeEntry, w io.Writer) error {
		t, err := util.NewTemplate(format)
		if err != nil {
			return err
		}

		l := len(timeEntries)
		for i := 0; i < l; i++ {
			if err := t.Execute(w, struct {
				dto.TimeEntry
				First bool
				Last  bool
			}{
				TimeEntry: timeEntries[i],
				First:     i == 0,
				Last:      i == (l - 1),
			}); err != nil {
				return err
			}
		}
		return nil
	}
}
