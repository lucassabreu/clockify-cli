package edit

import (
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/time-entry"
	"github.com/lucassabreu/clockify-cli/pkg/timeentryhlp"
	"github.com/spf13/cobra"
)

// NewCmdEdit represents the edit command
func NewCmdEdit(f cmdutil.Factory) *cobra.Command {
	of := util.OutputFlags{TimeFormat: output.TimeFormatSimple}
	cmd := &cobra.Command{
		Use: "edit [" +
			timeentryhlp.AliasCurrent + "|" + timeentryhlp.AliasLast +
			"|<time-entry-id>]",
		Aliases:   []string{"update"},
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{timeentryhlp.AliasLast, timeentryhlp.AliasCurrent},
		Short: `Edit a time entry, use id "` + timeentryhlp.AliasCurrent +
			`" to apply to time entry in progress`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
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

			tei, err := timeentryhlp.GetTimeEntry(
				c,
				w,
				userID,
				args[0],
			)
			if err != nil {
				return err
			}

			dc := util.NewDescriptionCompleter(f)

			if tei, err = util.ManageEntry(
				tei,
				util.FillTimeEntryWithFlags(cmd.Flags()),
				util.GetAllowNameForIDsFn(f.Config(), c),
				util.GetPropsInteractiveFn(c, dc, f.Config()),
				util.GetDatesInteractiveFn(f.Config()),
				util.GetValidateTimeEntryFn(f),
			); err != nil {
				return err
			}

			if tei, err = c.UpdateTimeEntry(api.UpdateTimeEntryParam{
				Workspace:   tei.WorkspaceID,
				TimeEntryID: tei.ID,
				Description: tei.Description,
				Start:       tei.TimeInterval.Start,
				End:         tei.TimeInterval.End,
				Billable:    tei.Billable,
				ProjectID:   tei.ProjectID,
				TaskID:      tei.TaskID,
				TagIDs:      tei.TagIDs,
			}); err != nil {
				return err
			}

			return util.PrintTimeEntryImpl(tei, f, cmd.OutOrStdout(), of)
		},
	}

	util.AddTimeEntryFlags(cmd, f, &of)

	cmd.Flags().StringP("when", "s", "",
		"when the entry should be started"+util.TimeFormatExamples)
	cmd.Flags().StringP("when-to-close", "e", "",
		"when the entry should be closed (same formats as `when`)")

	return cmd
}
