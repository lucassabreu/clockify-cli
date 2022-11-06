package edit

import (
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/time-entry"
	"github.com/lucassabreu/clockify-cli/pkg/timeentryhlp"
	"github.com/spf13/cobra"
)

// NewCmdEdit represents the edit command
func NewCmdEdit(
	f cmdutil.Factory,
	report func(dto.TimeEntryImpl, io.Writer, util.OutputFlags) error,
) *cobra.Command {
	of := util.OutputFlags{TimeFormat: output.TimeFormatSimple}
	va := cmdcompl.ValidArgsSlide{
		timeentryhlp.AliasCurrent, timeentryhlp.AliasLast}
	cmd := &cobra.Command{
		Use: "edit { <time-entry-id> | " + va.IntoUseOptions() +
			" | ^n }",
		Aliases: []string{"update"},
		Args: cobra.MatchAll(
			cmdutil.RequiredNamedArgs("time entry id"),
			cobra.ExactArgs(1),
		),
		ValidArgs: va.IntoValidArgs(),
		Short:     `Edit a time entry`,
		Long: heredoc.Docf(`
			Edit a time entry.
			Only the inputs sent thought flags will be changed, any other properties will remain the same.

			%s
			%s
			%s
			%s
			%s
		`,
			util.HelpTimeEntriesAliasForEdit,
			util.HelpInteractiveByDefault,
			util.HelpDateTimeFormats,
			util.HelpNamesForIds,
			util.HelpMoreInfoAboutPrinting,
		),
		Example: heredoc.Docf(`
			# starting a time entry
			$ %[1]s in --project cli --tag dev -d "Adding docs to edit" --task "edit" --md
			ID: %[2]s62ae4b304ebb4f143c931d50%[2]s  
			Billable: %[2]syes%[2]s  
			Locked: %[2]sno%[2]s  
			Project: Clockify Cli (%[2]s621948458cb9606d934ebb1c%[2]s)  
			Task: Edit Command (%[2]s62ae4af04ebb4f143c931d2e%[2]s)  
			Interval: %[2]s2022-06-18 22:01:16%[2]s until %[2]snow%[2]s  
			Description:
			> Adding docs to edit

			Tags:
			 * Development (%[2]s62ae28b72518aa18da2acb49%[2]s)

			# changing the description on the running time entry
			$ %[1]s edit current -d "Adding examples to edit" --md
			ID: %[2]s62ae4b304ebb4f143c931d50%[2]s  
			Billable: %[2]syes%[2]s  
			Locked: %[2]sno%[2]s  
			Project: Clockify Cli (%[2]s621948458cb9606d934ebb1c%[2]s)  
			Task: Edit Command (%[2]s62ae4af04ebb4f143c931d2e%[2]s)  
			Interval: %[2]s2022-06-18 22:01:16%[2]s until %[2]snow%[2]s  
			Description:
			> Adding examples to edit

			Tags:
			 * Development (%[2]s62ae28b72518aa18da2acb49%[2]s)

			# change the description, task, and tags
			$ %[1]s edit -d "Adding examples to edit" -T pair --task "in command" --md
			ID: %[2]s62ae4b304ebb4f143c931d50%[2]s  
			Billable: %[2]syes%[2]s
			Locked: %[2]sno%[2]s
			Project: Clockify Cli (%[2]s621948458cb9606d934ebb1c%[2]s)
			Task: In Command (%[2]s62ae29e62518aa18da2acd14%[2]s)
			Interval: %[2]s2022-06-18 22:13:14%[2]s until %[2]snow%[2]s
			Description:
			> Adding examples to edit

			Tags:
			 * Pair Programming (%[2]s621948708cb9606d934ebba7%[2]s)
		`, "clockify-cli", "`"),
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

			if tei, err = util.Do(
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

			if report != nil {
				return report(tei, cmd.OutOrStdout(), of)
			}

			return util.PrintTimeEntryImpl(tei, f, cmd.OutOrStdout(), of)
		},
	}

	util.AddTimeEntryFlags(cmd, f, &of)

	cmd.Flags().StringP("when", "s", "",
		"when the entry should be started")
	cmd.Flags().StringP("when-to-close", "e", "",
		"when the entry should be closed (same formats as `when`)")

	return cmd
}
