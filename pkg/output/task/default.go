package task

import (
	"io"
	"os"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/term"
)

// TaskPrint will print more details
func TaskPrint(ts []dto.Task, w io.Writer) error {
	tw := tablewriter.NewWriter(w)
	tw.SetHeader([]string{"ID", "Name", "Status"})

	lines := make([][]string, len(ts))
	for i := 0; i < len(ts); i++ {
		lines[i] = []string{
			ts[i].ID,
			ts[i].Name,
			string(ts[i].Status),
		}
	}

	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		tw.SetColWidth(width / 3)
	}
	tw.AppendBulk(lines)
	tw.Render()

	return nil
}
