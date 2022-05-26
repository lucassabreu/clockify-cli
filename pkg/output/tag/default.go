package tag

import (
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/olekukonko/tablewriter"
)

// TagPrint will print more details
func TagPrint(ts []dto.Tag, w io.Writer) error {
	tw := tablewriter.NewWriter(w)
	tw.SetHeader([]string{"ID", "Name"})

	lines := make([][]string, len(ts))
	for i := 0; i < len(ts); i++ {
		lines[i] = []string{
			ts[i].ID,
			ts[i].Name,
		}
	}

	tw.AppendBulk(lines)
	tw.Render()

	return nil
}
