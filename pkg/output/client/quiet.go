package output

import (
	"fmt"
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// ClientPrintQuietly will only print the IDs
func ClientPrintQuietly(cs []dto.Client, w io.Writer) error {
	for i := 0; i < len(cs); i++ {
		fmt.Fprintln(w, cs[i].ID)
	}

	return nil
}
