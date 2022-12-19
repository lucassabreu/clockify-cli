package util

import (
	"time"

	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcomplutil"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
)

// AddTimeEntryFlags will add the common flags needed to add/edit a time entry
func AddTimeEntryFlags(
	cmd *cobra.Command, f cmdutil.Factory, of *OutputFlags,
) {
	cmd.Flags().BoolP("billable", "b", false,
		"this time entry is billable")
	cmd.Flags().BoolP("not-billable", "n", false,
		"this time entry is not billable")
	cmd.Flags().String("task", "", "add a task to the entry")
	_ = cmdcompl.AddSuggestionsToFlag(cmd, "task",
		cmdcomplutil.NewTaskAutoComplete(f, true))

	cmd.Flags().StringSliceP("tag", "T", []string{}, "add tags to the entry (can be used multiple times)")
	_ = cmdcompl.AddSuggestionsToFlag(cmd, "tag",
		cmdcomplutil.NewTagAutoComplete(f))

	cmd.Flags().BoolP("allow-incomplete", "A", false,
		"allow creation of incomplete time entries to be edited later")

	cmd.Flags().StringP("project", "p", "", "project to use for time entry")
	_ = cmdcompl.AddSuggestionsToFlag(cmd, "project",
		cmdcomplutil.NewProjectAutoComplete(f))

	cmd.Flags().StringP("description", "d", "", "time entry description")
	_ = cmdcompl.AddSuggestionsToFlag(cmd, "description",
		newDescriptionAutoComplete(f),
	)

	AddPrintTimeEntriesFlags(cmd, of)

	// deprecations
	cmd.Flags().StringSlice("tags", []string{}, "add tags to the entry")
	_ = cmdcompl.AddSuggestionsToFlag(cmd, "tags",
		cmdcomplutil.NewTagAutoComplete(f))
	_ = cmd.Flags().MarkDeprecated("tags", "use tag instead")
}

// AddTimeEntryDateFlags adds the default start and end flags
func AddTimeEntryDateFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("when", "s", time.Now().Format(timehlp.FullTimeFormat),
		"when the entry should be started, "+
			"if not informed will use current time")
	cmd.Flags().StringP("when-to-close", "e", "",
		"when the entry should be closed, if not informed will let it open "+
			"(same formats as when)")
}
