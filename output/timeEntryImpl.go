package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/template"
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/crypto/ssh/terminal"
)

// TimeEntryImplJSONPrint will print as JSON
func TimeEntryImplJSONPrint(t *dto.TimeEntryImpl, w io.Writer) error {
	return json.NewEncoder(w).Encode(t)
}

// TimeEntryImplPrint will print more details
func TimeEntryImplPrint(t *dto.TimeEntryImpl, w io.Writer) error {
	tw := tablewriter.NewWriter(w)
	tw.SetHeader([]string{"Start", "End", "Dur", "Project", "Description"})

	var wD, wP = 0, 0

	if wD < len(t.Description) {
		wD = len(t.Description)
	}

	if wP < len(t.ProjectID) {
		wP = len(t.ProjectID)
	}

	if width, _, err := terminal.GetSize(int(os.Stdin.Fd())); err == nil {
		width = width - 30 - wP
		if width < wD {
			wD = width
		}
	}

	end := time.Now()
	if t.TimeInterval.End != nil {
		end = *t.TimeInterval.End
	}

	end = end.Round(time.Second)

	tw.SetColWidth(wD)
	tw.AppendBulk([][]string{
		{
			t.TimeInterval.Start.In(time.Local).Format("2006-01-02 15:04:05"),
			end.In(time.Local).Format("2006-01-02 15:04:05"),
			fmt.Sprintf("%-8v", end.Sub(t.TimeInterval.Start)),
			t.ProjectID,
			t.Description,
		},
	})
	tw.Render()

	return nil
}

// TimeEntryImplPrintWithTemplate will print each time entry using the format string
func TimeEntryImplPrintWithTemplate(format string) func(*dto.TimeEntryImpl, io.Writer) error {
	return func(tei *dto.TimeEntryImpl, w io.Writer) error {
		t, err := template.New("tmpl").Parse(format)
		if err != nil {
			return err
		}

		if err := t.Execute(w, tei); err != nil {
			return err
		}
		fmt.Fprintln(w)
		return nil
	}
}
