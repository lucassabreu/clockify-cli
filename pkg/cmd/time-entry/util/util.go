package util

import (
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/time-entry"
	"github.com/spf13/cobra"
)

// PrintTimeEntries will print out a time entries using parameters and flags
func PrintTimeEntryImpl(
	tei dto.TimeEntryImpl,
	f cmdutil.Factory,
	cmd *cobra.Command,
	timeFormat string,
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

	return PrintTimeEntry(fte, cmd, timeFormat, f.Config())
}

// PrintTimeEntries will print out a time entries using parameters and flags
func PrintTimeEntry(te *dto.TimeEntry,
	cmd *cobra.Command, timeFormat string, config cmdutil.Config) error {
	ts := make([]dto.TimeEntry, 0)
	if te != nil {
		ts = append(ts, *te)
	}

	b := config.GetBool(cmdutil.CONF_SHOW_TOTAL_DURATION)
	config.SetBool(cmdutil.CONF_SHOW_TOTAL_DURATION, false)

	err := PrintTimeEntries(ts, cmd, timeFormat, config)

	config.SetBool(cmdutil.CONF_SHOW_TOTAL_DURATION, b)

	return err
}

// PrintTimeEntries will print out a list of time entries using parameters and
// flags
func PrintTimeEntries(
	tes []dto.TimeEntry,
	cmd *cobra.Command, timeFormat string, config cmdutil.Config,
) error {
	out := cmd.OutOrStdout()
	if b, _ := cmd.Flags().GetBool("md"); b {
		return output.TimeEntriesMarkdownPrint(tes, out)
	}

	if asJSON, _ := cmd.Flags().GetBool("json"); asJSON {
		return output.TimeEntriesJSONPrint(tes, out)
	}

	if asCSV, _ := cmd.Flags().GetBool("csv"); asCSV {
		return output.TimeEntriesCSVPrint(tes, out)
	}

	if format, _ := cmd.Flags().GetString("format"); format != "" {
		return output.TimeEntriesPrintWithTemplate(format)(tes, out)
	}

	if asQuiet, _ := cmd.Flags().GetBool("quiet"); asQuiet {
		return output.TimeEntriesPrintQuietly(tes, out)
	}

	if b, _ := cmd.Flags().GetBool("duration-float"); b {
		return output.TimeEntriesTotalDurationOnlyAsFloat(tes, out)
	}

	if b, _ := cmd.Flags().GetBool("duration-formatted"); b {
		return output.TimeEntriesTotalDurationOnlyFormatted(tes, out)
	}

	opts := []output.TimeEntryOutputOpt{output.WithTimeFormat(timeFormat)}

	if config.GetBool(cmdutil.CONF_SHOW_TASKS) {
		opts = append(opts, output.WithShowTasks())
	}

	if config.GetBool(cmdutil.CONF_SHOW_TOTAL_DURATION) {
		opts = append(opts, output.WithTotalDuration())
	}

	return output.TimeEntriesPrint(opts...)(tes, out)
}
