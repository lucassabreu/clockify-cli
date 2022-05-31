package show

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timeentryhlp"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
)

// NewCmdShow represents the show command
func NewCmdShow(f cmdutil.Factory) *cobra.Command {
	of := util.OutputFlags{TimeFormat: timehlp.FullTimeFormat}
	cmd := &cobra.Command{
		Use: "show [" + timeentryhlp.AliasCurrent + "|" +
			timeentryhlp.AliasLast + "|<time-entry-id>|^n]",
		ValidArgs: []string{timeentryhlp.AliasCurrent, timeentryhlp.AliasLast},
		Args:      cobra.MaximumNArgs(1),
		Short:     "Show detailed information about one time entry.",
		Long: heredoc.Doc(`
			Show detailed information about one time entry.
			Shows current one by default
		`),
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
