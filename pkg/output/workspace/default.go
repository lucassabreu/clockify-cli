package workspace

import (
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/olekukonko/tablewriter"
)

// WorkspacePrint will print more details
func WorkspacePrint(
	wDefault string) func(ws []dto.Workspace, w io.Writer) error {
	return func(ws []dto.Workspace, w io.Writer) error {
		tw := tablewriter.NewWriter(w)
		tw.SetHeader([]string{"ID", "Name", "Image"})

		lines := make([][]string, len(ws))
		for i := 0; i < len(ws); i++ {
			lines[i] = []string{
				ws[i].ID,
				ws[i].Name,
				ws[i].ImageURL,
			}
			if wDefault == ws[i].ID {
				lines[i][1] = lines[i][1] + " (default)"
			}
		}

		tw.AppendBulk(lines)
		tw.Render()

		return nil
	}
}
