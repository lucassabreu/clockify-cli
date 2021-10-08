package output

import (
	"embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/ui"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/term"
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

func colorToTermColor(hex string) []int {
	if len(hex) == 0 {
		return []int{}
	}

	fi, _ := os.Stdout.Stat()
	if fi.Mode()&os.ModeCharDevice == 0 {
		return []int{}
	}

	if c, err := ui.HEX(hex[1:]); err == nil {
		return append(
			[]int{38, 2},
			c.Values()...,
		)
	}

	return []int{}
}

//go:embed resources
var res embed.FS

// TimeEntriesMarkdownPrint will print time entries in "markdown blocks"
func TimeEntriesMarkdownPrint(tes []dto.TimeEntry, w io.Writer) error {
	b, err := res.ReadFile("resources/timeEntry.gotmpl.md")
	if err != nil {
		return err
	}

	return TimeEntriesPrintWithTemplate(string(b))(tes, w)
}

// TimeEntryMarkdownPrint will print time entries in "markdown blocks"
func TimeEntryMarkdownPrint(te *dto.TimeEntry, w io.Writer) error {
	if te == nil {
		return nil
	}

	return TimeEntriesMarkdownPrint([]dto.TimeEntry{*te}, w)
}

// TimeEntriesPrint will print more details
func TimeEntriesPrint(showTasks bool, timeFormat ...string) func([]dto.TimeEntry, io.Writer) error {
	format := TIME_FORMAT_SIMPLE
	if len(timeFormat) > 0 {
		format = timeFormat[0]
	}

	return func(timeEntries []dto.TimeEntry, w io.Writer) error {
		tw := tablewriter.NewWriter(w)
		header := []string{"ID", "Start", "End", "Dur",
			"Project", "Description", "Tags"}
		if showTasks {
			header = append(
				header[:5],
				header[5:]...,
			)
			header[5] = "Task"
		}

		tw.SetHeader(header)
		tw.SetRowLine(true)
		if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			tw.SetColWidth(width / 3)
		}

		colors := make([]tablewriter.Colors, len(header))
		for _, t := range timeEntries {
			end := time.Now()
			if t.TimeInterval.End != nil {
				end = *t.TimeInterval.End
			}

			projectName := ""
			colors[4] = []int{}
			if t.Project != nil {
				colors[4] = colorToTermColor(t.Project.Color)
				projectName = t.Project.Name
			}

			line := []string{
				t.ID,
				t.TimeInterval.Start.In(time.Local).Format(format),
				end.In(time.Local).Format(format),
				durationToString(end.Sub(t.TimeInterval.Start)),
				projectName,
				t.Description,
				strings.Join(tagsToStringSlice(t.Tags), ", "),
			}

			if showTasks {
				line = append(line[:5], line[5:]...)
				line[5] = ""
				if t.Task != nil {
					line[5] = fmt.Sprintf("%s (%s)", t.Task.Name, t.Task.ID)
				}
			}

			tw.Rich(line, colors)
		}

		tw.Render()

		return nil
	}
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
		"task.id",
		"task.name",
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

		if te.Task == nil {
			t := dto.Task{}
			te.Task = &t
		}

		arr := []string{
			te.ID,
			te.Description,
			p.ID,
			p.Name,
			te.Task.ID,
			te.Task.Name,
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

var funcMap = template.FuncMap{
	"formatDateTime": func(t time.Time) string {
		return t.Format(TIME_FORMAT_FULL)
	},
}

// TimeEntriesPrintWithTemplate will print each time entry using the format
// string
func TimeEntriesPrintWithTemplate(
	format string,
) func([]dto.TimeEntry, io.Writer) error {
	return func(timeEntries []dto.TimeEntry, w io.Writer) error {
		t, err := template.New("tmpl").Funcs(funcMap).Parse(format)
		if err != nil {
			return err
		}

		for i, te := range timeEntries {
			if err := t.Execute(w, struct {
				dto.TimeEntry
				First bool
				Last  bool
			}{
				TimeEntry: te,
				First:     i == 0,
				Last:      i == (len(timeEntries) - 1),
			}); err != nil {
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

// TimeEntryCSVPrint will print as CSV
func TimeEntryCSVPrint(t *dto.TimeEntry, w io.Writer) error {
	entries := []dto.TimeEntry{}

	if t != nil {
		entries = append(entries, *t)
	}

	return TimeEntriesCSVPrint(entries, w)
}

// TimeEntryPrintQuietly will only print the IDs
func TimeEntryPrintQuietly(timeEntry *dto.TimeEntry, w io.Writer) error {
	fmt.Fprintln(w, timeEntry.ID)
	return nil
}

// TimeEntryPrint will print more details
func TimeEntryPrint(
	showTasks bool, timeFormat ...string,
) func(*dto.TimeEntry, io.Writer) error {
	return func(timeEntry *dto.TimeEntry, w io.Writer) error {
		entries := []dto.TimeEntry{}

		if timeEntry != nil {
			entries = append(entries, *timeEntry)
		}

		return TimeEntriesPrint(showTasks, timeFormat...)(entries, w)
	}
}

// TimeEntryPrintWithTemplate will print each time entry using the format string
func TimeEntryPrintWithTemplate(
	format string,
) func(*dto.TimeEntry, io.Writer) error {
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
	return fmt.Sprintf("%d:%02d:%02d",
		int64(d.Hours()), int64(d.Minutes())%60, int64(d.Seconds())%60)
}
