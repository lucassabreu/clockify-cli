package timeentry

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/output/util"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/term"
)

func sumTimeEntriesDuration(ts []dto.TimeEntry) time.Duration {
	s := time.Duration(0)
	for i := 0; i < len(ts); i++ {
		end := time.Now()
		if ts[i].TimeInterval.End != nil {
			end = *ts[i].TimeInterval.End
		}

		d := end.Sub(ts[i].TimeInterval.Start)
		s = s + d
	}
	return s
}

const (
	TimeFormatFull   = "2006-01-02 15:04:05"
	TimeFormatSimple = "15:04:05"
)

// TimeEntryOutputOptions sets how the "table" format should print the time
// entries
type TimeEntryOutputOptions struct {
	ShowTasks         bool
	ShowClients       bool
	ShowTotalDuration bool
	TimeFormat        string
}

// NewTimeEntryOutputOptions creates a default TimeEntryOutputOptions
func NewTimeEntryOutputOptions() TimeEntryOutputOptions {
	return TimeEntryOutputOptions{
		TimeFormat:        TimeFormatSimple,
		ShowTasks:         false,
		ShowClients:       false,
		ShowTotalDuration: false,
	}
}

// WithTimeFormat sets the date-time output format
func (teo TimeEntryOutputOptions) WithTimeFormat(
	format string) TimeEntryOutputOptions {
	teo.TimeFormat = format
	return teo

}

// WithShowTasks shows a new column with the task of the time entry
func (teo TimeEntryOutputOptions) WithShowTasks() TimeEntryOutputOptions {
	teo.ShowTasks = true
	return teo
}

// WithShowCliens shows a new column with the client of the time entry
func (teo TimeEntryOutputOptions) WithShowClients() TimeEntryOutputOptions {
	teo.ShowClients = true
	return teo
}

// WithTotalDuration shows a footer with the sum of the durations of the time
// entries
func (teo TimeEntryOutputOptions) WithTotalDuration() TimeEntryOutputOptions {
	teo.ShowTotalDuration = true
	return teo
}

// TimeEntriesPrint will print more details
func TimeEntriesPrint(
	options TimeEntryOutputOptions) func([]dto.TimeEntry, io.Writer) error {
	return func(timeEntries []dto.TimeEntry, w io.Writer) error {
		tw := tablewriter.NewWriter(w)
		projectColumn := 4
		header := []string{"ID", "Start", "End", "Dur", "Project"}

		if options.ShowClients {
			header = append(header, "Client")
		}

		if options.ShowTasks {
			header = append(header, "Task")
		}

		header = append(header, "Description", "Tags")

		tw.SetHeader(header)
		tw.SetRowLine(true)
		if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			if options.ShowClients || options.ShowTasks {
				tw.SetColWidth(width / 4)
			} else {
				tw.SetColWidth(width / 3)
			}
		}

		colors := make([]tablewriter.Colors, len(header))
		for i := 0; i < len(timeEntries); i++ {
			t := timeEntries[i]
			end := time.Now()
			if t.TimeInterval.End != nil {
				end = *t.TimeInterval.End
			}

			projectName := ""
			colors[projectColumn] = []int{}
			if t.Project != nil {
				colors[projectColumn] = util.ColorToTermColor(t.Project.Color)
				projectName = t.Project.Name
			}

			line := []string{
				t.ID,
				t.TimeInterval.Start.In(time.Local).Format(options.TimeFormat),
				end.In(time.Local).Format(options.TimeFormat),
				durationToString(end.Sub(t.TimeInterval.Start)),
				projectName,
			}

			if options.ShowClients {
				client := ""
				if t.Project.ClientName != "" {
					colors[len(line)] = colors[projectColumn]
					client = t.Project.ClientName
				}
				line = append(line, client)
			}

			if options.ShowTasks {
				task := ""
				if t.Task != nil {
					task = fmt.Sprintf("%s (%s)", t.Task.Name, t.Task.ID)
				}
				line = append(line, task)
			}

			line = append(
				line,
				t.Description,
				strings.Join(tagsToStringSlice(t.Tags), "\n"),
			)

			tw.Rich(line, colors)
		}

		if options.ShowTotalDuration {
			line := make([]string, len(header))
			line[0] = "TOTAL"
			line[3] = durationToString(sumTimeEntriesDuration(timeEntries))
			tw.Append(line)
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

func durationToString(d time.Duration) string {
	return dto.Duration{Duration: d}.HumanString()
}
