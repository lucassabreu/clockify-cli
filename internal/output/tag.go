package output

import (
	"fmt"
	"io"
	"text/template"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/olekukonko/tablewriter"
)

// TagPrintQuietly will only print the IDs
func TagPrintQuietly(ws []dto.Tag, w io.Writer) error {
	for _, wk := range ws {
		fmt.Fprintln(w, wk.ID)
	}

	return nil
}

// TagPrint will print more details
func TagPrint(ws []dto.Tag, w io.Writer) error {
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

// TagPrintWithTemplate will print each worspace using the format string
func TagPrintWithTemplate(format string) func([]dto.Tag, io.Writer) error {
	return func(ws []dto.Tag, w io.Writer) error {
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
