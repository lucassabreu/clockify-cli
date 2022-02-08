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

// ClientPrintQuietly will only print the IDs
func ClientPrintQuietly(cs []dto.Client, w io.Writer) error {
	for _, c := range cs {
		fmt.Fprintln(w, c.ID)
	}

	return nil
}

// ClientPrint will print more details
func ClientPrint(cs []dto.Client, w io.Writer) error {
	tw := tablewriter.NewWriter(w)
	tw.SetHeader([]string{"ID", "Name", "Archived"})

	yesNo := map[bool]string{
		true:  "YES",
		false: "NO",
	}

	lines := make([][]string, len(cs))
	for i, c := range cs {
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

// ClientPrintWithTemplate will print each client using the format string
func ClientPrintWithTemplate(format string) func([]dto.Client, io.Writer) error {
	return func(ws []dto.Client, w io.Writer) error {
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

// ClientsJSONPrint will print as JSON
func ClientsJSONPrint(t []dto.Client, w io.Writer) error {
	return json.NewEncoder(w).Encode(t)
}

// ClientsCSVPrint will print as CSV
func ClientsCSVPrint(clients []dto.Client, out io.Writer) error {
	w := csv.NewWriter(out)

	err := w.Write([]string{
		"id",
		"name",
		"archived",
	})

	if err != nil {
		return err
	}

	for _, c := range clients {
		arr := []string{
			c.ID,
			c.Name,
			fmt.Sprintf("%v", c.Archived),
		}

		err := w.Write(arr)

		if err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}
