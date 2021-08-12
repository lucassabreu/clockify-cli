package output

import (
	"encoding/json"
	"fmt"
	"io"
	"text/template"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/olekukonko/tablewriter"
)

// UserPrintQuietly will only print the IDs
func UserPrintQuietly(users []dto.User, w io.Writer) error {
	for _, u := range users {
		fmt.Fprintln(w, u.ID)
	}

	return nil
}

func UserJSONPrint(u dto.User, w io.Writer) error {
	return json.NewEncoder(w).Encode(u)
}

// UserPrint will print more details
func UserPrint(users []dto.User, w io.Writer) error {
	tw := tablewriter.NewWriter(w)
	tw.SetHeader([]string{"ID", "Name", "Email", "Status"})

	lines := make([][]string, len(users))
	for i, u := range users {
		lines[i] = []string{
			u.ID,
			u.Name,
			u.Email,
			string(u.Status),
		}
	}

	tw.AppendBulk(lines)
	tw.Render()

	return nil
}

// UserPrintWithTemplate will print each worspace using the format string
func UserPrintWithTemplate(format string) func([]dto.User, io.Writer) error {
	return func(users []dto.User, w io.Writer) error {
		t, err := template.New("tmpl").Parse(format)
		if err != nil {
			return err
		}

		for _, i := range users {
			if err := t.Execute(w, i); err != nil {
				return err
			}
			fmt.Fprintln(w)
		}
		return nil
	}
}
