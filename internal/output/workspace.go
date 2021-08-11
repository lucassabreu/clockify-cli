package output

import (
	"fmt"
	"io"
	"text/template"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/olekukonko/tablewriter"
)

// WorkspacePrintQuietly will only print the IDs
func WorkspacePrintQuietly(ws []dto.Workspace, w io.Writer) error {
	for _, wk := range ws {
		fmt.Fprintln(w, wk.ID)
	}

	return nil
}

// WorkspacePrint will print more details
func WorkspacePrint(wDefault string) func(ws []dto.Workspace, w io.Writer) error {
	return func(ws []dto.Workspace, w io.Writer) error {
		tw := tablewriter.NewWriter(w)
		tw.SetHeader([]string{"ID", "Name", "Image"})

		lines := make([][]string, len(ws))
		for i, w := range ws {
			name := w.Name
			if wDefault == w.ID {
				name = name + " (default)"
			}
			lines[i] = []string{
				w.ID,
				name,
				w.ImageURL,
			}
		}

		tw.AppendBulk(lines)
		tw.Render()

		return nil
	}
}

// WorkspacePrintWithTemplate will print each worspace using the format string
func WorkspacePrintWithTemplate(format string) func([]dto.Workspace, io.Writer) error {
	return func(ws []dto.Workspace, w io.Writer) error {
		t, err := template.New("tmpl").Parse(format)
		if err != nil {
			return err
		}

		for _, i := range ws {
			if err := t.Execute(w, i); err != nil {
				return err
			}
			fmt.Fprintln(w)
		}
		return nil
	}
}
