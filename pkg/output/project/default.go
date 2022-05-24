package project

import (
	"fmt"
	"io"
	"os"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/term"
)

// ProjectPrint will print more details
func ProjectPrint(ps []dto.Project, w io.Writer) error {
	tw := tablewriter.NewWriter(w)
	tw.SetHeader([]string{"ID", "Name", "Client"})

	lines := make([][]string, len(ps))
	for i := 0; i < len(ps); i++ {
		w := ps[i]
		client := ""
		if w.ClientID != "" {
			client = fmt.Sprintf("%s (%s)", w.ClientName, w.ClientID)
		}

		lines[i] = []string{
			w.ID,
			w.Name,
			client,
		}
	}

	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		tw.SetColWidth(width / 3)
	}
	tw.AppendBulk(lines)
	tw.Render()

	return nil
}
