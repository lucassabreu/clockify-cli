package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/term"
)

// TaskPrintQuietly will only print the IDs
func TaskPrintQuietly(cs []dto.Task, w io.Writer) error {
	for _, c := range cs {
		fmt.Fprintln(w, c.ID)
	}

	return nil
}

// TaskPrint will print more details
func TaskPrint(ts []dto.Task, w io.Writer) error {
	tw := tablewriter.NewWriter(w)
	tw.SetHeader([]string{"ID", "Name", "Status"})

	lines := make([][]string, len(ts))
	for i, t := range ts {
		lines[i] = []string{
			t.ID,
			t.Name,
			string(t.Status),
		}
	}

	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		tw.SetColWidth(width / 3)
	}
	tw.AppendBulk(lines)
	tw.Render()

	return nil
}

// TaskPrintWithTemplate will print each client using the format string
func TaskPrintWithTemplate(format string) func([]dto.Task, io.Writer) error {
	return func(ws []dto.Task, w io.Writer) error {
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

// TasksJSONPrint will print as JSON
func TasksJSONPrint(t []dto.Task, w io.Writer) error {
	return json.NewEncoder(w).Encode(t)
}

// TasksCSVPrint will print as CSV
func TasksCSVPrint(tasks []dto.Task, out io.Writer) error {
	w := csv.NewWriter(out)

	err := w.Write([]string{
		"id",
		"name",
		"status",
	})

	if err != nil {
		return err
	}

	for _, t := range tasks {
		arr := []string{
			t.ID,
			t.Name,
			string(t.Status),
		}

		err := w.Write(arr)

		if err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}
