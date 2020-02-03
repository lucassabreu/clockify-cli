package reports

import (
	"encoding/csv"
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

// TimeEntriesJSONPrint will print as JSON
func TimeEntriesJSONPrint(t []dto.TimeEntry, w io.Writer) error {
	return json.NewEncoder(w).Encode(t)
}

// TimeEntriesPrintQuietly will only print the IDs
func TimeEntriesPrintQuietly(timeEntries []dto.TimeEntry, w io.Writer) error {
	for _, u := range timeEntries {
		fmt.Fprintln(w, u.ID)
	}

	return nil
}

// TimeEntriesPrint will print more details
func TimeEntriesPrint(timeEntries []dto.TimeEntry, w io.Writer) error {
	tw := tablewriter.NewWriter(w)
	tw.SetHeader([]string{"ID", "Start", "End", "Dur", "Project", "Description"})

	lines := make([][]string, len(timeEntries))

	var wD, wP = 0, 0

	for i, t := range timeEntries {
		if wD < len(t.Description) {
			wD = len(t.Description)
		}

		if t.Project != nil && wP < len(t.Project.Name) {
			wP = len(t.Project.Name)
		}

		end := time.Now()
		if t.TimeInterval.End != nil {
			end = *t.TimeInterval.End
		}

		projectName := ""
		if t.Project != nil {
			projectName = t.Project.Name
		}
		lines[i] = []string{
			t.ID,
			t.TimeInterval.Start.In(time.Local).Format("15:04:05"),
			end.In(time.Local).Format("15:04:05"),
			fmt.Sprintf("%-8v", end.Sub(t.TimeInterval.Start)),
			projectName,
			t.Description,
		}
	}

	if width, _, err := terminal.GetSize(int(os.Stdin.Fd())); err == nil {
		width = width - 70 - wP
		if width < wD {
			wD = width
		}
	}

	tw.SetColWidth(wD)
	tw.AppendBulk(lines)
	tw.Render()

	return nil
}

// TimeEntriesCSVPrint will print each time entry using the format string
func TimeEntriesCSVPrint(timeEntries []dto.TimeEntry, out io.Writer) error {
	w := csv.NewWriter(out)

	err := w.Write([]string{
		"id",
		"description",
		"project.id",
		"project.name",
		"start",
		"end",
		"duration",
		"user.id",
		"user.email",
		"user.name",
		"tags...",
	})

	if err != nil {
		return err
	}

	format := func(t *time.Time) string {
		if t == nil {
			return ""
		}
		return t.In(time.Local).Format("2006-01-02 15:04:05")
	}

	tags := func(tags []dto.Tag) []string {
		s := make([]string, len(tags))

		for i, t := range tags {
			s[i] = fmt.Sprintf("%s (%s)", t.Name, t.ID)
		}

		return s
	}

	for _, te := range timeEntries {
		var p dto.Project
		if te.Project != nil {
			p = *te.Project
		}

		arr := []string{
			te.ID,
			te.Description,
			p.ID,
			p.Name,
			format(&te.TimeInterval.Start),
			format(te.TimeInterval.End),
			fmt.Sprintf("%-8v", te.TimeInterval.End.Sub(te.TimeInterval.Start)),
			te.User.ID,
			te.User.Email,
			te.User.Name,
		}

		err := w.Write(append(arr, tags(te.Tags)...))

		if err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}

// TimeEntriesPrintWithTemplate will print each time entry using the format string
func TimeEntriesPrintWithTemplate(format string) func([]dto.TimeEntry, io.Writer) error {
	return func(timeEntries []dto.TimeEntry, w io.Writer) error {
		t, err := template.New("tmpl").Parse(format)
		if err != nil {
			return err
		}

		for _, i := range timeEntries {
			if err := t.Execute(w, i); err != nil {
				return err
			}
			fmt.Fprintln(w)
		}
		return nil
	}
}

// TimeEntryJSONPrint will print as JSON
func TimeEntryJSONPrint(t *dto.TimeEntry, w io.Writer) error {
	return json.NewEncoder(w).Encode(t)
}

// TimeEntryPrintQuietly will only print the IDs
func TimeEntryPrintQuietly(timeEntry *dto.TimeEntry, w io.Writer) error {
	fmt.Fprintln(w, timeEntry.ID)
	return nil
}

// TimeEntryPrint will print more details
func TimeEntryPrint(timeEntry *dto.TimeEntry, w io.Writer) error {
	entries := []dto.TimeEntry{}

	if timeEntry != nil {
		entries = append(entries, *timeEntry)
	}

	return TimeEntriesPrint(entries, w)
}

// TimeEntryPrintWithTemplate will print each time entry using the format string
func TimeEntryPrintWithTemplate(format string) func(*dto.TimeEntry, io.Writer) error {
	fn := TimeEntriesPrintWithTemplate(format)
	return func(timeEntry *dto.TimeEntry, w io.Writer) error {
		entries := []dto.TimeEntry{}

		if timeEntry != nil {
			entries = append(entries, *timeEntry)
		}

		return fn(entries, w)
	}
}
