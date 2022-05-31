package editmultiple

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/time-entry"
	"github.com/lucassabreu/clockify-cli/pkg/timeentryhlp"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/spf13/cobra"
)

// NewCmdEditMultiple represents the editMultiple command
func NewCmdEditMultiple(f cmdutil.Factory) *cobra.Command {
	of := util.OutputFlags{TimeFormat: output.TimeFormatSimple}
	cmd := &cobra.Command{
		Use: "edit-multiple [" +
			timeentryhlp.AliasCurrent + "|" + timeentryhlp.AliasLast +
			"|<time-entry-id>]...",
		Aliases: []string{
			"update-multiple", "multi-edit",
			"multi-update", "mult-edit", "mult-update",
		},
		Args:      cobra.MinimumNArgs(2),
		ValidArgs: []string{timeentryhlp.AliasLast, timeentryhlp.AliasCurrent},
		Short: `Edit multiple time entries at once, use id "` +
			timeentryhlp.AliasCurrent + `"/"` + timeentryhlp.AliasLast +
			`" to apply to time entry in progress.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var w, u string

			if w, err = f.GetWorkspaceID(); err != nil {
				return err
			}

			if u, err = f.GetUserID(); err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			teis := make([]dto.TimeEntryImpl, len(args))
			for i := range args {
				if teis[i], err = timeentryhlp.GetTimeEntry(
					c, w, u, args[i]); err != nil {
					return err
				}
			}

			tei := teis[0]
			editFn := func(tei dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
				return c.UpdateTimeEntry(api.UpdateTimeEntryParam{
					Workspace:   tei.WorkspaceID,
					TimeEntryID: tei.ID,
					Description: tei.Description,
					Start:       tei.TimeInterval.Start,
					End:         tei.TimeInterval.End,
					Billable:    tei.Billable,
					ProjectID:   tei.ProjectID,
					TaskID:      tei.TaskID,
					TagIDs:      tei.TagIDs,
				})
			}

			fn := func(input dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
				var err error
				for i, tei := range teis {
					input.TimeInterval = tei.TimeInterval
					input.ID = tei.ID

					if tei, err = editFn(input); err != nil {
						return input, err
					}

					teis[i] = tei
				}

				return input, err
			}

			if !f.Config().IsInteractive() {
				fn = func(input dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
					c := cmd.Flags().Changed
					for i, tei := range teis {
						if c("project") {
							tei.ProjectID = input.ProjectID
						}

						if c("description") {
							tei.Description = input.Description
						}

						if c("task") {
							tei.TaskID = input.TaskID
						}

						if c("tag") || c("tags") {
							tei.TagIDs = input.TagIDs
						}

						if c("not-billable") {
							tei.Billable = input.Billable
						}

						teis[i] = tei
						if _, err = editFn(tei); err != nil {
							return tei, err
						}
					}
					return input, nil
				}
			}

			dc := util.NewDescriptionCompleter(f)

			if _, err = util.ManageEntry(
				tei,
				util.FillTimeEntryWithFlags(cmd.Flags()),
				util.GetAllowNameForIDsFn(f.Config(), c),
				util.GetPropsInteractiveFn(c, dc, f.Config()),
				util.GetValidateTimeEntryFn(f),
				fn,
			); err != nil {
				return err
			}

			tes := make([]dto.TimeEntry, len(teis))
			var t *dto.TimeEntry
			for i, tei := range teis {
				t, err = c.GetHydratedTimeEntry(api.GetTimeEntryParam{
					TimeEntryID: tei.ID,
					Workspace:   tei.WorkspaceID,
				})

				if err != nil {
					return err
				}
				tes[i] = *t
			}

			return util.PrintTimeEntries(tes,
				cmd.OutOrStdout(), f.Config(), of)
		},
	}

	cmd.Long = cmd.Short + heredoc.Doc(`
		When multiple IDs are informed the default values on interactive mode will be the values of the first time entry informed.
		When using interactive mode all entries will end with the same properties except for Start and End, if you wanna edit only some properties, than use the flags without interactive mode.
		Start and end fields can't be mass-edited.
	`)

	util.AddTimeEntryFlags(cmd, f, &of)
	util.AddPrintMultipleTimeEntriesFlags(cmd)

	return cmd
}
