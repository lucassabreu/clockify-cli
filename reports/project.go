package reports

import (
	"fmt"
	"io"
	"text/template"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/olekukonko/tablewriter"
)

// ProjectPrintQuietly will only print the IDs
func ProjectPrintQuietly(ws []dto.Project, w io.Writer) error {
	for _, wk := range ws {
		fmt.Fprintln(w, wk.ID)
	}

	return nil
}

// ProjectPrint will print more details
func ProjectPrint(ws []dto.Project, w io.Writer) error {
	tw := tablewriter.NewWriter(w)
	tw.SetHeader([]string{"ID", "Name"})

	lines := make([][]string, len(ws))
	for i, w := range ws {
		lines[i] = []string{
			w.ID,
			w.Name,
		}
	}

	tw.AppendBulk(lines)
	tw.Render()

	return nil
}

// ProjectPrintWithTemplate will print each worspace using the format string
func ProjectPrintWithTemplate(format string) func([]dto.Project, io.Writer) error {
	return func(ws []dto.Project, w io.Writer) error {
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
