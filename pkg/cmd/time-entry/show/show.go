package show

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timeentryhlp"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
)

// NewCmdShow represents the show command
func NewCmdShow(f cmdutil.Factory) *cobra.Command {
	of := util.OutputFlags{TimeFormat: timehlp.FullTimeFormat}
	va := cmdcompl.ValidArgsSlide{
		timeentryhlp.AliasCurrent, timeentryhlp.AliasLast}
	cmd := &cobra.Command{
		Use: "show [ <time-entry-id> | " + va.IntoUseOptions() +
			" | ^n ]",
		ValidArgs: va.IntoValidArgs(),
		Args:      cobra.MaximumNArgs(1),
		Short:     "Show information about one time entry.",
		Long: heredoc.Docf(`
			Show information about one time entry.

			If no time entry ID is informed it shows the running it exists.

			To show the last ended time entry you can use "%s" for it, for the one before that you can use "^2", for the previous "^3" and so on.

			%s
		`,
			timeentryhlp.AliasLast,
			util.HelpMoreInfoAboutPrinting,
		),
		Example: heredoc.Docf(`
			# trying to show running time entry, when there is none
			$ %[1]s
			looking for running time entry: time entry was not found

			# show the last time entry (ended)
			$ %[1]s last -q
			62af70d849445270d7c09fbd

			# show the time entry before the last one
			$ %[1]s ^2 -q
			62af668b49445270d7c092e4
		`, "clockify-cli show"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			userID, err := f.GetUserID()
			if err != nil {
				return err
			}

			w, err := f.GetWorkspaceID()
			if err != nil {
				return err
			}

			id := timeentryhlp.AliasCurrent
			if len(args) > 0 {
				id = args[0]
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			tei, err := timeentryhlp.GetTimeEntry(c, w, userID, id)
			if err != nil {
				return err
			}

			return util.PrintTimeEntryImpl(tei, f, cmd.OutOrStdout(), of)
		},
	}

	util.AddPrintTimeEntriesFlags(cmd, &of)
	_ = cmd.MarkFlagRequired("workspace")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}
