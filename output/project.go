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
	"golang.org/x/crypto/ssh/terminal"
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
	tw.SetHeader([]string{"ID", "Name", "Client"})

	lines := make([][]string, len(ws))
	for i, w := range ws {

		client := ""
		if len(w.ClientID) != 0 {
			client = fmt.Sprintf("%s (%s)", w.ClientName, w.ClientID)
		}

		lines[i] = []string{
			w.ID,
			w.Name,
			client,
		}
	}

	if width, _, err := terminal.GetSize(int(os.Stdin.Fd())); err == nil {
		tw.SetColWidth(width / 3)
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

// ProjectsJSONPrint will print as JSON
func ProjectsJSONPrint(t []dto.Project, w io.Writer) error {
	return json.NewEncoder(w).Encode(t)
}

// ProjectsCSVPrint will print each time entry using the format string
func ProjectsCSVPrint(projects []dto.Project, out io.Writer) error {
	w := csv.NewWriter(out)

	err := w.Write([]string{
		"id",
		"name",
		"client.id",
		"client.name",
	})

	if err != nil {
		return err
	}

	for _, p := range projects {
		arr := []string{
			p.ID,
			p.Name,
			p.ClientID,
			p.ClientName,
		}

		err := w.Write(arr)

		if err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}
