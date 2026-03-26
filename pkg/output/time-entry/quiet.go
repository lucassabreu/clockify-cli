package timeentry

import (
	"fmt"
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// TimeEntriesPrintQuietly will only print the IDs
func TimeEntriesPrintQuietly(timeEntries []dto.TimeEntry, w io.Writer) error {
	for i := 0; i < len(timeEntries); i++ {
		if _, err := fmt.Fprintln(w, timeEntries[i].ID); err != nil {
			return err
		}
	}

	return nil
}
