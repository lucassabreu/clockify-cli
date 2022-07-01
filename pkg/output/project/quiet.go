package project

import (
	"fmt"
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// ProjectPrintQuietly will only print the IDs
func ProjectPrintQuietly(ps []dto.Project, w io.Writer) error {
	for i := 0; i < len(ps); i++ {
		fmt.Fprintln(w, ps[i].ID)
	}

	return nil
}
