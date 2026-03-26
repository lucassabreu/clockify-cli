package tag

import (
	"fmt"
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// TagPrintQuietly will only print the IDs
func TagPrintQuietly(ts []dto.Tag, w io.Writer) error {
	for i := 0; i < len(ts); i++ {
		if _, err := fmt.Fprintln(w, ts[i].ID); err != nil {
			return err
		}
	}

	return nil
}
