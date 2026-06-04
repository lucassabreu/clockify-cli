package edit

import (
	"errors"
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
			" | ^n }...",
		Aliases: []string{
			"update",
			"update-multiple", "multi-edit",
			"multi-update", "mult-edit", "mult-update",
		},
		Args: cobra.MatchAll(
			cmdutil.RequiredNamedArgs("time entry id"),
			cobra.MinimumNArgs(1),
		),
		ValidArgs: va.IntoValidArgs(),
		Short:     `Edit one or more time entries`,
		Long: heredoc.Docf(`
			Edit one or more time entries.

			When editing a single time entry, you can use --when and --when-to-close to change when it started or ended.
			Only the inputs sent thought flags will be changed, any other properties will remain the same.

			When editing multiple time entries, you can change all properties except for when they start or end,
			as different time entries will have different start and end times.

			Except on interactive mode where the values informed, even if not changed will be applied to all entries
			(except for Start and End time).

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

			# just to help show the data
			$ export F="{{.ID}} :: {{ .Description }}
			  When: {{ fdt .TimeInterval.Start }} util {{ ft (.TimeInterval.End | now) }}
			  Task: {{ .Task.Name }} ({{ .Project.Name}})
			  Tags: {{ .Tags }}
			"

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

			if len(args) > 1 {
				if cmd.Flags().Changed("when") || cmd.Flags().Changed("when-to-close") {
					return errors.New("--when and --when-to-close can only be used when editing a single time entry")
				}
			}

			teis := make([]util.TimeEntryDTO, len(args))
			for i := range args {
				t, err := timeentryhlp.GetTimeEntry(c, w, userID, args[i])
				if err != nil {
					return err
				}
				teis[i] = util.TimeEntryImplToDTO(t)
			}

			dc := util.NewDescriptionCompleter(f)

			if len(args) == 1 {
				te := teis[0]
				if te, err = util.Do(
					te,
					util.FillTimeEntryWithFlags(cmd.Flags()),
					util.GetAllowNameForIDsFn(f.Config(), c),
					util.GetPropsInteractiveFn(dc, f),
					util.GetDatesInteractiveFn(f),
					util.GetValidateTimeEntryFn(f),
				); err != nil {
					return err
				}

				tei, err := c.UpdateTimeEntry(api.UpdateTimeEntryParam{
					Workspace:   te.Workspace,
					TimeEntryID: te.ID,
					Description: te.Description,
					Start:       te.Start,
					End:         te.End,
					Billable:    *te.Billable,
					ProjectID:   te.ProjectID,
					TaskID:      te.TaskID,
					TagIDs:      te.TagIDs,
				})
				if err != nil {
					return err
				}

				return report(tei, cmd.OutOrStdout(), of)
			}

			tei := teis[0]
			editFn := func(tei util.TimeEntryDTO) (util.TimeEntryDTO, error) {
				t, err := c.UpdateTimeEntry(api.UpdateTimeEntryParam{
					Workspace:   tei.Workspace,
					TimeEntryID: tei.ID,
					Description: tei.Description,
					Start:       tei.Start,
					End:         tei.End,
					Billable:    *tei.Billable,
					ProjectID:   tei.ProjectID,
					TaskID:      tei.TaskID,
					TagIDs:      tei.TagIDs,
				})

				return util.TimeEntryImplToDTO(t), err
			}

			fn := func(input util.TimeEntryDTO) (util.TimeEntryDTO, error) {
				var err error
				for i, tei := range teis {
					input.Start = tei.Start
					input.End = tei.End
					input.ID = tei.ID

					if tei, err = editFn(input); err != nil {
						return input, err
					}

					teis[i] = tei
				}

				return input, err
			}

			if !f.Config().IsInteractive() {
				fn = func(input util.TimeEntryDTO) (util.TimeEntryDTO, error) {
					changed := cmd.Flags().Changed

					if changed("task") && !changed("project") {
						projectID := teis[0].ProjectID
						for _, te := range teis[1:] {
							if te.ProjectID != projectID {
								return input, errors.New("you are changing the task of the time entries, but not the project and some of them are not in the same project, please also set --project")
							}
						}
					}

					for i, tei := range teis {
						if changed("project") {
							if tei.ProjectID != input.ProjectID {
								tei.TaskID = ""
							}
							tei.ProjectID = input.ProjectID
						}

						if changed("description") {
							tei.Description = input.Description
						}

						if changed("task") {
							tei.TaskID = input.TaskID
						}

						if changed("tag") || changed("tags") {
							tei.TagIDs = input.TagIDs
						}

						if changed("not-billable") {
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

			if _, err = util.Do(
				tei,
				util.FillTimeEntryWithFlags(cmd.Flags()),
				util.GetAllowNameForIDsFn(f.Config(), c),
				util.GetPropsInteractiveFn(dc, f),
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
					Workspace:   tei.Workspace,
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

	cmd.Flags().StringP("when", "s", "",
		"when the entry should be started")
	cmd.Flags().StringP("when-to-close", "e", "",
		"when the entry should be closed (same formats as `when`)")

	return cmd
}
