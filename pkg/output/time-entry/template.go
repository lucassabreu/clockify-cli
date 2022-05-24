package timeentry

import (
	"fmt"
	"io"
	"text/template"
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

var funcMap = template.FuncMap{
	"formatDateTime": func(t time.Time) string {
		return t.Format(TimeFormatFull)
	},
}

// TimeEntriesPrintWithTemplate will print each time entry using the format
// string
func TimeEntriesPrintWithTemplate(
	format string,
) func([]dto.TimeEntry, io.Writer) error {
	return func(timeEntries []dto.TimeEntry, w io.Writer) error {
		t, err := template.New("tmpl").Funcs(funcMap).Parse(format)
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
			fmt.Fprintln(w)
		}
		return nil
	}
}
