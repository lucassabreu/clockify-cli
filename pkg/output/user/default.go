package user

import (
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/olekukonko/tablewriter"
)

// UserPrint will print more details
func UserPrint(users []dto.User, w io.Writer) error {
	tw := tablewriter.NewWriter(w)
	tw.SetHeader([]string{"ID", "Name", "Email", "Status", "TimeZone"})

	lines := make([][]string, len(users))
	for i := 0; i < len(users); i++ {
		lines[i] = []string{
			users[i].ID,
			users[i].Name,
			users[i].Email,
			string(users[i].Status),
			users[i].Settings.TimeZone,
		}
	}

	tw.AppendBulk(lines)
	tw.Render()

	return nil
}
