package task

import (
	"fmt"
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// TaskPrintQuietly will only print the IDs
func TaskPrintQuietly(ts []dto.Task, w io.Writer) error {
	for i := 0; i < len(ts); i++ {
		fmt.Fprintln(w, ts[i].ID)
	}

	return nil
}
