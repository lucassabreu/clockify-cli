package task

import (
	"encoding/csv"
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// TasksCSVPrint will print as CSV
func TasksCSVPrint(ts []dto.Task, out io.Writer) error {
	w := csv.NewWriter(out)

	if err := w.Write([]string{
		"id",
		"name",
		"status",
	}); err != nil {
		return err
	}

	for i := 0; i < len(ts); i++ {
		if err := w.Write([]string{
			ts[i].ID,
			ts[i].Name,
			string(ts[i].Status),
		}); err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}
