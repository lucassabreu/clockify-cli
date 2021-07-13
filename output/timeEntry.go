package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
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

const (
	TIME_FORMAT_FULL   = "2006-01-02 15:04:05"
	TIME_FORMAT_SIMPLE = "15:04:05"
)

// TimeEntriesPrintWithTimeFormat will print more details
func TimeEntriesPrintWithTimeFormat(format string) func([]dto.TimeEntry, io.Writer) error {
	return func(timeEntries []dto.TimeEntry, w io.Writer) error {
		tw := tablewriter.NewWriter(w)
		tw.SetHeader([]string{"ID", "Start", "End", "Dur", "Project", "Description", "Tags"})

		lines := make([][]string, len(timeEntries))

		for i, t := range timeEntries {
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
				t.TimeInterval.Start.In(time.Local).Format(format),
				end.In(time.Local).Format(format),
				durationToString(end.Sub(t.TimeInterval.Start)),
				projectName,
				t.Description,
				strings.Join(tagsToStringSlice(t.Tags), ", "),
			}
		}

		if width, _, err := terminal.GetSize(int(os.Stdin.Fd())); err == nil {
			tw.SetColWidth(width / 3)
		}

		tw.SetRowLine(true)
		tw.AppendBulk(lines)
		tw.Render()

		return nil
	}
}

// TimeEntriesPrint will print more details
func TimeEntriesPrint(timeEntries []dto.TimeEntry, w io.Writer) error {
	return TimeEntriesPrintWithTimeFormat(TIME_FORMAT_SIMPLE)(timeEntries, w)
}

func tagsToStringSlice(tags []dto.Tag) []string {
	s := make([]string, len(tags))

	for i, t := range tags {
		s[i] = fmt.Sprintf("%s (%s)", t.Name, t.ID)
	}

	return s
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

	for _, te := range timeEntries {
		var p dto.Project
		if te.Project != nil {
			p = *te.Project
		}

		end := time.Now()
		if te.TimeInterval.End != nil {
			end = *te.TimeInterval.End
		}

		if te.User == nil {
			u := dto.User{}
			te.User = &u
		}

		arr := []string{
			te.ID,
			te.Description,
			p.ID,
			p.Name,
			format(&te.TimeInterval.Start),
			format(te.TimeInterval.End),
			durationToString(end.Sub(te.TimeInterval.Start)),
			te.User.ID,
			te.User.Email,
			te.User.Name,
		}

		err := w.Write(append(arr, tagsToStringSlice(te.Tags)...))

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

func durationToString(d time.Duration) string {
	return fmt.Sprintf("%d:%02d:%02d", int64(d.Hours()), int64(d.Minutes())%60, int64(d.Seconds())%60)
}
