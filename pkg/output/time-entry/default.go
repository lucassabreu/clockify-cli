package timeentry

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/ui"
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

// TimeEntryOptions sets how the "table" format should print the time entries
type TimeEntryOutputOptions struct {
	ShowTasks         bool
	ShowTotalDuration bool
	TimeFormat        string
}

// WithTimeFormat sets the date-time output format
func WithTimeFormat(format string) TimeEntryOutputOpt {
	return func(teo *TimeEntryOutputOptions) error {
		teo.TimeFormat = format
		return nil
	}
}

// WithShowTasks shows a new column with the task of the time entry
func WithShowTasks() TimeEntryOutputOpt {
	return func(teoo *TimeEntryOutputOptions) error {
		teoo.ShowTasks = true
		return nil
	}
}

// WithDurationTotal shows a footer with the sum of the durations of the time
// entries
func WithTotalDuration() TimeEntryOutputOpt {
	return func(teoo *TimeEntryOutputOptions) error {
		teoo.ShowTotalDuration = true
		return nil
	}
}

// TimeEntryOutputOpt allows the setting of TimeEntryOutputOptions values
type TimeEntryOutputOpt func(*TimeEntryOutputOptions) error

// TimeEntriesPrint will print more details
func TimeEntriesPrint(opts ...TimeEntryOutputOpt) func([]dto.TimeEntry, io.Writer) error {
	options := &TimeEntryOutputOptions{
		TimeFormat:        TimeFormatSimple,
		ShowTasks:         false,
		ShowTotalDuration: false,
	}

	for _, o := range opts {
		err := o(options)
		if err != nil {
			return func(te []dto.TimeEntry, w io.Writer) error { return err }
		}
	}

	return func(timeEntries []dto.TimeEntry, w io.Writer) error {
		tw := tablewriter.NewWriter(w)
		taskColumn := 6
		projectColumn := 4
		header := []string{"ID", "Start", "End", "Dur",
			"Project", "Description", "Tags"}
		if options.ShowTasks {
			header = append(
				header[:taskColumn],
				header[taskColumn-1:]...,
			)
			header[taskColumn] = "Task"
		}

		tw.SetHeader(header)
		tw.SetRowLine(true)
		if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			tw.SetColWidth(width / 3)
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
				colors[projectColumn] = colorToTermColor(t.Project.Color)
				projectName = t.Project.Name
			}

			line := []string{
				t.ID,
				t.TimeInterval.Start.In(time.Local).Format(options.TimeFormat),
				end.In(time.Local).Format(options.TimeFormat),
				durationToString(end.Sub(t.TimeInterval.Start)),
				projectName,
				t.Description,
				strings.Join(tagsToStringSlice(t.Tags), "\n"),
			}

			if options.ShowTasks {
				line = append(line[:taskColumn], line[taskColumn-1:]...)
				line[taskColumn] = ""
				if t.Task != nil {
					line[taskColumn] = fmt.Sprintf("%s (%s)", t.Task.Name, t.Task.ID)
				}
			}

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
	p := ""
	if d < 0 {
		p = "-"
		d = d * -1
	}

	return p + fmt.Sprintf("%d:%02d:%02d",
		int64(d.Hours()), int64(d.Minutes())%60, int64(d.Seconds())%60)
}
