package client

import (
	"io"
	"os"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/term"
)

// ClientPrint will print more details
func ClientPrint(cs []dto.Client, w io.Writer) error {
	tw := tablewriter.NewWriter(w)
	tw.SetHeader([]string{"ID", "Name", "Archived"})

	yesNo := map[bool]string{
		true:  "YES",
		false: "NO",
	}

	lines := make([][]string, len(cs))
	for i := 0; i < len(cs); i++ {
		c := cs[i]
		lines[i] = []string{
			c.ID,
			c.Name,
			yesNo[c.Archived],
		}
	}

	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		tw.SetColWidth(width / 3)
	}
	tw.AppendBulk(lines)
	tw.Render()

	return nil
}
