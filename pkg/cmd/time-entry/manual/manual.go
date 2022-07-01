package manual

import (
	"fmt"
	"time"

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

// NewCmdManual represents the manual command
func NewCmdManual(f cmdutil.Factory) *cobra.Command {
	of := util.OutputFlags{TimeFormat: output.TimeFormatSimple}
	cmd := &cobra.Command{
		Use:   "manual [<project-id>] [<start>] [<end>] [<description>]",
		Short: "Create a new complete time entry",
		Long: heredoc.Doc(`
			Create a new complete time entry with start and end.

			This command will not stop running time entries.

			The rules defined in the workspace and project will be checked before creating it.
		`) + "\n" +
			util.HelpTimeEntryNowIfNotSet +
			"The same applies to end time (`--when-to-close`).\n\n" +
			util.HelpInteractiveByDefault + "\n" +
			util.HelpTimeInputOnTimeEntry + "\n" +
			util.HelpNamesForIds + "\n" +
			util.HelpMoreInfoAboutStarting + "\n" +
			util.HelpMoreInfoAboutPrinting,
		Args: cobra.MaximumNArgs(4),
		ValidArgsFunction: cmdcompl.CombineSuggestionsToArgs(
			cmdcomplutil.NewProjectAutoComplete(f)),
		RunE: func(cmd *cobra.Command, args []string) error {
			var whenToCloseDate time.Time
			var err error
			tei := dto.TimeEntryImpl{}

			if tei.WorkspaceID, err = f.GetWorkspaceID(); err != nil {
				return err
			}

			if tei.UserID, err = f.GetUserID(); err != nil {
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
				tei.TimeInterval.Start, err = timehlp.ConvertToTime(args[1])
				if err != nil {
					return fmt.Errorf(
						"Fail to convert when to start: %w", err)
				}
			}

			if len(args) > 2 {
				whenToCloseDate, err = timehlp.ConvertToTime(args[2])
				if err != nil {
					return fmt.Errorf(
						"Fail to convert when to end: %w", err)
				}
				tei.TimeInterval.End = &whenToCloseDate
			}

			if len(args) > 3 {
				tei.Description = args[3]
			}

			dc := util.NewDescriptionCompleter(f)

			if tei, err = util.Do(
				tei,
				util.FillTimeEntryWithFlags(cmd.Flags()),
				func(tei dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
					if tei.TimeInterval.End != nil {
						return tei, nil
					}

					now, _ := timehlp.ConvertToTime(timehlp.NowTimeFormat)
					tei.TimeInterval.End = &now
					return tei, nil
				},
				util.GetAllowNameForIDsFn(f.Config(), c),
				util.GetPropsInteractiveFn(c, dc, f.Config()),
				util.GetDatesInteractiveFn(f.Config()),
				util.ValidateClosingTimeEntry(f),
				util.CreateTimeEntryFn(c),
			); err != nil {
				return err
			}

			return util.PrintTimeEntryImpl(tei,
				f, cmd.OutOrStdout(), of)
		},
	}

	util.AddTimeEntryFlags(cmd, f, &of)
	util.AddTimeEntryDateFlags(cmd)

	return cmd
}
