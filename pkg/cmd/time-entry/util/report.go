package util

import (
	"io"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/time-entry"
	"github.com/spf13/cobra"
)

// OutputFlags sets how to print out a list of time entries
type OutputFlags struct {
	Format            string
	CSV               bool
	JSON              bool
	Quiet             bool
	Markdown          bool
	DurationFormatted bool
	DurationFloat     bool
	TimeFormat        string
}

func (of OutputFlags) Check() error {
	return cmdutil.XorFlag(map[string]bool{
		"format":             of.Format != "",
		"json":               of.JSON,
		"csv":                of.CSV,
		"quiet":              of.Quiet,
		"md":                 of.Markdown,
		"duration-float":     of.DurationFloat,
		"duration-formatted": of.DurationFormatted,
	})
}

// AddPrintMultipleTimeEntriesFlags add flags to print multiple time entries
func AddPrintMultipleTimeEntriesFlags(cmd *cobra.Command) {
	cmd.Flags().BoolP("with-totals", "S", false,
		"add a totals line at the end")
}

// AddPrintTimeEntriesFlags add flags common to time entry print
func AddPrintTimeEntriesFlags(cmd *cobra.Command, of *OutputFlags) {
	cmd.Flags().StringVarP(&of.Format, "format", "f", "",
		"golang text/template format to be applied on each time entry")
	cmd.Flags().String("tz", "Local",
		"time zone to be used on the time entries can be "+
			"'Local' to use the systems timezone, UTC "+
			"or valid TZ identifier from the IANA TZ database "+
			"https://en.wikipedia.org/wiki/List_of_tz_database_time_zones")
	cmd.Flags().BoolVarP(&of.JSON, "json", "j", false, "print as JSON")
	cmd.Flags().BoolVarP(&of.CSV, "csv", "v", false, "print as CSV")
	cmd.Flags().BoolVarP(&of.Quiet, "quiet", "q", false, "print only ID")
	cmd.Flags().BoolVarP(&of.Markdown, "md", "m", false, "print as Markdown")
	cmd.Flags().BoolVarP(&of.DurationFormatted, "duration-formatted", "D", false,
		"prints only the sum of duration formatted")
	cmd.Flags().BoolVarP(&of.DurationFloat, "duration-float", "F", false,
		`prints only the sum of duration as a "float hour"`)
}

// PrintTimeEntryImpl will print out a time entries using parameters and flags
func PrintTimeEntryImpl(
	tei dto.TimeEntryImpl,
	f cmdutil.Factory,
	out io.Writer,
	of OutputFlags,
) error {
	c, err := f.Client()
	if err != nil {
		return err
	}

	fte, err := c.GetHydratedTimeEntry(api.GetTimeEntryParam{
		Workspace:   tei.WorkspaceID,
		TimeEntryID: tei.ID,
	})

	if err != nil {
		return err
	}

	return PrintTimeEntry(fte, out, f.Config(), of)
}

// PrintTimeEntry will print out a time entries using parameters and flags
func PrintTimeEntry(
	te *dto.TimeEntry, out io.Writer, config cmdutil.Config, of OutputFlags,
) error {
	ts := make([]dto.TimeEntry, 0)
	if te != nil {
		ts = append(ts, *te)
	}

	b := config.GetBool(cmdutil.CONF_SHOW_TOTAL_DURATION)
	config.SetBool(cmdutil.CONF_SHOW_TOTAL_DURATION, false)

	err := PrintTimeEntries(ts, out, config, of)

	config.SetBool(cmdutil.CONF_SHOW_TOTAL_DURATION, b)

	return err
}

func updateTimeZone(tes []dto.TimeEntry, config cmdutil.Config) []dto.TimeEntry {
	loc := config.TimeZone()
	if loc == time.Local {
		return tes
	}

	for i := range tes {
		tes[i].TimeInterval.Start = tes[i].TimeInterval.Start.In(loc)
		if tes[i].TimeInterval.End != nil {
			end := tes[i].TimeInterval.End.In(loc)
			tes[i].TimeInterval.End = &end
		}
	}
	return tes
}

// PrintTimeEntries will print out a list of time entries using parameters and
// flags
func PrintTimeEntries(
	tes []dto.TimeEntry, out io.Writer, config cmdutil.Config, of OutputFlags,
) error {
	tes = updateTimeZone(tes, config)
	switch {
	case of.Markdown:
		return output.TimeEntriesMarkdownPrint(tes, out)
	case of.JSON:
		return output.TimeEntriesJSONPrint(tes, out)
	case of.CSV:
		return output.TimeEntriesCSVPrint(tes, out)
	case of.Format != "":
		return output.TimeEntriesPrintWithTemplate(of.Format)(tes, out)
	case of.Quiet:
		return output.TimeEntriesPrintQuietly(tes, out)
	case of.DurationFloat:
		return output.TimeEntriesTotalDurationOnlyAsFloat(
			tes, out, config.Language())
	case of.DurationFormatted:
		return output.TimeEntriesTotalDurationOnlyFormatted(tes, out)
	default:
		opts := output.NewTimeEntryOutputOptions().
			WithTimeFormat(of.TimeFormat)

		if config.GetBool(cmdutil.CONF_SHOW_TASKS) {
			opts = opts.WithShowTasks()
		}

		if config.GetBool(cmdutil.CONF_SHOW_CLIENT) {
			opts = opts.WithShowClients()
		}

		if config.GetBool(cmdutil.CONF_SHOW_TOTAL_DURATION) {
			opts = opts.WithTotalDuration()
		}

		return output.TimeEntriesPrint(opts)(tes, out)
	}
}
