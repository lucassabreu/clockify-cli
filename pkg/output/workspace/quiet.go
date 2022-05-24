package workspace

import (
	"fmt"
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// WorkspacePrintQuietly will only print the IDs
func WorkspacePrintQuietly(ws []dto.Workspace, w io.Writer) error {
	for i := 0; i < len(ws); i++ {
		fmt.Fprintln(w, ws[i].ID)
	}

	return nil
}
