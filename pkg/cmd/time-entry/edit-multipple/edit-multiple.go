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
		Use: "edit-multiple { <time-entry-id> | " +
			timeentryhlp.AliasCurrent + " | " + timeentryhlp.AliasLast +
			" }...",
		Aliases: []string{
			"update-multiple", "multi-edit",
			"multi-update", "mult-edit", "mult-update",
		},
		Args: cobra.MatchAll(
			cmdutil.RequiredNamedArgs("time entry id"),
			cobra.MinimumNArgs(2),
		),
		ValidArgs: []string{timeentryhlp.AliasLast, timeentryhlp.AliasCurrent},
		Short:     `Edit multiple time entries at once`,
		Long: heredoc.Docf(`
			Edit multiple time entries at once.

			This command does not allow to edit when the time entries start or ended, because different time entries will have different start and end times.

			Except on interactive mode where the values informed, even if not changed will be applied to all entries (except for Start and End time).
			If you wanna edit only some properties, than use the flags without interactive mode, only the input sent thought the flags will be changed.

			%s
			%s
			%s
			%s
		`,
			util.HelpTimeEntriesAliasForEdit,
			util.HelpInteractiveByDefault,
			util.HelpNamesForIds,
			util.HelpMoreInfoAboutPrinting,
		),
		Example: heredoc.Docf(`
			# just to help show the data
			$ export F="{{.ID}} :: {{ .Description }}
			  When: {{ fdt .TimeInterval.Start }} util {{ ft (.TimeInterval.End | now) }}
			  Task: {{ .Task.Name }} ({{ .Project.Name}})
			  Tags: {{ .Tags }}
			"

			$ %[1]s report --format "$F"
			62af667c4ebb4f143c9482bb :: Edit multiple entries
			  When: 2022-06-19 18:10:01 util 18:10:15
			  Task: Edit Command (Clockify Cli)
			  Tags: [Development (62ae28b72518aa18da2acb49)]

			62af668b49445270d7c092e4 :: Adding examples
			  When: 2022-06-19 18:10:15 util 18:29:32
			  Task: Edit Command (Clockify Cli)
			  Tags: [Development (62ae28b72518aa18da2acb49)]

			62af6b0f4ebb4f143c94880e :: More examples
			  When: 2022-06-19 18:29:32 util 18:38:12
			  Task: Edit Command (Clockify Cli)
			  Tags: [Development (62ae28b72518aa18da2acb49)]

			# change all to use other task
			$ %[1]s edit-multiple -i=0 -f "$F" current last ^2 --task multiple
			62af6b0f4ebb4f143c94880e :: More examples
			  When: 2022-06-19 18:29:32 util 18:43:04
			  Task: Edit Multiple Command (Clockify Cli)
			  Tags: [Development (62ae28b72518aa18da2acb49)]
			62af668b49445270d7c092e4 :: Adding examples
			  When: 2022-06-19 18:10:15 util 18:29:32
			  Task: Edit Multiple Command (Clockify Cli)
			  Tags: [Development (62ae28b72518aa18da2acb49)]
			62af668b49445270d7c092e4 :: Adding examples
			  When: 2022-06-19 18:10:15 util 18:29:32
			  Task: Edit Multiple Command (Clockify Cli)
			  Tags: [Development (62ae28b72518aa18da2acb49)]
		`, "clockify-cli"),
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

			if _, err = util.Do(
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

	util.AddTimeEntryFlags(cmd, f, &of)
	util.AddPrintMultipleTimeEntriesFlags(cmd)

	return cmd
}
