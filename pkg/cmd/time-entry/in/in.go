package in

import (
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcomplutil"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/time-entry"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"

	"github.com/spf13/cobra"
)

// NewCmdIn represents the in command
func NewCmdIn(
	f cmdutil.Factory,
	report func(dto.TimeEntryImpl, io.Writer, util.OutputFlags) error,
) *cobra.Command {
	of := util.OutputFlags{TimeFormat: output.TimeFormatSimple}
	cmd := &cobra.Command{
		Use:   "in [<project-id>] [<description>]",
		Short: "Create a new Clockify time entry ",
		Long: heredoc.Doc(`
			Create a new Clockify time entry

			Running time entry will be stopped using the start time of this new entry.
		`) + "\n" +
			util.HelpTimeEntryNowIfNotSet + "\n" +
			util.HelpInteractiveByDefault + "\n" +
			util.HelpTimeInputOnTimeEntry + "\n" +
			util.HelpNamesForIds + "\n" +
			util.HelpValidateIncomplete + "\n" +
			util.HelpMoreInfoAboutPrinting,
		Args: cobra.MaximumNArgs(2),
		ValidArgsFunction: cmdcompl.CombineSuggestionsToArgs(
			cmdcomplutil.NewProjectAutoComplete(f, f.Config())),
		Aliases: []string{"start"},
		Example: heredoc.Docf(`
			# start a timer with project and description, starting now
			$ %[1]s -i=0 "Clockify CLI" "Documenting in command"
			+--------------------------+----------+----------+---------+--------------+------------------------+------+
			|            ID            |  START   |   END    |   DUR   |   PROJECT    |      DESCRIPTION       | TAGS |
			+--------------------------+----------+----------+---------+--------------+------------------------+------+
			| 62ae2744c22de9759e73d038 | 13:28:01 | 13:28:04 | 0:00:03 | Clockify Cli | Documenting in command |      |
			+--------------------------+----------+----------+---------+--------------+------------------------+------+

			# start a timer with description, starting at 14:00
			$ %[1]s -i=0 -d "Documenting in command" -s "14:00"
			+--------------------------+----------+----------+---------+---------+------------------------+------+
			|            ID            |  START   |   END    |   DUR   | PROJECT |      DESCRIPTION       | TAGS |
			+--------------------------+----------+----------+---------+---------+------------------------+------+
			| 62ae27cd49445270d7bf0333 | 14:00:00 | 14:30:21 | 0:30:21 |         | Documenting in command |      |
			+--------------------------+----------+----------+---------+---------+------------------------+------+

			# start a timer with description, project and tags, starting 10 min ago
			$ %[1]s -i=0 -p 621948458cb9606d934ebb1c -d "Documenting in command" -s -10m --tag dev
			+--------------------------+----------+----------+---------+--------------+------------------------+--------------------------------+
			|            ID            |  START   |   END    |   DUR   |   PROJECT    |      DESCRIPTION       |              TAGS              |
			+--------------------------+----------+----------+---------+--------------+------------------------+--------------------------------+
			| 62ae29104ebb4f143c92f458 | 14:25:41 | 14:35:44 | 0:10:03 | Clockify Cli | Documenting in command | Development                    |
			|                          |          |          |         |              |                        | (62ae28b72518aa18da2acb49)     |
			+--------------------------+----------+----------+---------+--------------+------------------------+--------------------------------+

			# start a timer with description, project and task, starting at 10 min, but only showing its ID
			$ %[1]s -i=0 -p 621948458cb9606d934ebb1c -d "Documenting in command" -s -10m --task "in command"
			62ae29fdc22de9759e73d343

			# start a timer without description, with task and project
			$ %[1]s -i=0 -p 621948458cb9606d934ebb1c -s -10m --task "in command"
			62ae29fdc22de9759e73d343

			# start a timer interactively
			$ %[1]s -i
			? Choose your project: 621948458cb9606d934ebb1c - Clockify Cli      | Client: Myself (6202634a28782767054eec26)
			? Choose your task: 62ae29e62518aa18da2acd14 - In Command
			? Description: Adding more examples
			? Choose your tags: 62ae28b72518aa18da2acb49 - Development, 621948708cb9606d934ebba7 - Pair Programming
			? Start: now
			? End (leave it blank for empty):
			+--------------------------+----------+----------+---------+--------------+----------------------+-----------------------------------------+
			|            ID            |  START   |   END    |   DUR   |   PROJECT    |     DESCRIPTION      |                  TAGS                   |
			+--------------------------+----------+----------+---------+--------------+----------------------+-----------------------------------------+
			| 62ae37b84ebb4f143c930523 | 17:38:14 | 17:38:17 | 0:00:03 | Clockify Cli | Adding more examples | Pair Programming                        |
			|                          |          |          |         |              |                      | (621948708cb9606d934ebba7) Development  |
			|                          |          |          |         |              |                      | (62ae28b72518aa18da2acb49)              |
			+--------------------------+----------+----------+---------+--------------+----------------------+-----------------------------------------+
		`, "clockify-cli in"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			var err error
			tei := util.TimeEntryDTO{
				Start: timehlp.Now(),
			}

			if tei.Workspace, err = f.GetWorkspaceID(); err != nil {
				return err
			}

			if tei.UserID, err = f.GetUserID(); err != nil {
				return err
			}

			tei, err = util.FromDefaults(f)(tei)
			if err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			if len(args) > 0 {
				tei.ProjectID = args[0]
			}

			if len(args) > 1 {
				tei.Description = args[1]
			}

			dc := util.NewDescriptionCompleter(f)

			if tei, err = util.Do(
				tei,

				util.FillTimeEntryWithFlags(cmd.Flags()),
				util.ValidateClosingTimeEntry(f),
				util.GetAllowNameForIDsFn(f.Config(), c),
				util.GetPropsInteractiveFn(dc, f),
				util.GetDatesInteractiveFn(f),
				util.GetValidateTimeEntryFn(f),
				util.OutInProgressFn(c),
				util.CreateTimeEntryFn(c),
			); err != nil {
				return err
			}

			return report(
				util.TimeEntryDTOToImpl(tei), cmd.OutOrStdout(), of)
		},
	}

	util.AddTimeEntryFlags(cmd, f, &of)
	util.AddTimeEntryDateFlags(cmd)

	return cmd
}
