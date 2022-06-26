package clone

import (
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	timeentry "github.com/lucassabreu/clockify-cli/pkg/output/time-entry"
	"github.com/lucassabreu/clockify-cli/pkg/timeentryhlp"
	"github.com/spf13/cobra"
)

// NewCmdClone represents the clone command
func NewCmdClone(f cmdutil.Factory) *cobra.Command {
	of := util.OutputFlags{TimeFormat: timeentry.TimeFormatSimple}
	cmd := &cobra.Command{
		Use: "clone " +
			"{ <time-entry-id> | " + timeentryhlp.AliasLast + "| ^<n> }",
		Short: "Copy a time entry and starts it ",
		Long: heredoc.Docf(`
			Copy a time entry and starts it.

			Running time entry will be stopped using the start time of this new entry. If you don't want to stop them, use the flag %[1]s--no-closing%[1]s.

			If you want to clone the last (running) time entry you can use "%[2]s" instead of its ID.
			Also if you want to clone the one previous to it, you can use "^2", for the before that "^3" and so on.

			The rules defined in the workspace and project will be checked before creating it.
		`, "`", timeentryhlp.AliasLast) + "\n" +
			util.HelpTimeEntryNowIfNotSet +
			"The same applies to end time (`--when-to-close`).\n\n" +
			util.HelpInteractiveByDefault + "\n" +
			util.HelpTimeInputOnTimeEntry + "\n" +
			util.HelpNamesForIds + "\n" +
			util.HelpMoreInfoAboutStarting + "\n" +
			util.HelpMoreInfoAboutPrinting,
		Example: heredoc.Docf(`
			$ %[1]s in --project cli --tag dev -d "Adding docs to clone" --task "clone" --md
			ID: %[2]s62ae4b304ebb4f143c931d50%[2]s  
			Billable: %[2]syes%[2]s  
			Locked: %[2]sno%[2]s  
			Project: Clockify Cli (%[2]s621948458cb9606d934ebb1c%[2]s)  
			Task: Clone Command (%[2]s62ae4af04ebb4f143c931d2e%[2]s)  
			Interval: %[2]s2022-06-18 22:01:16%[2]s until %[2]snow%[2]s  
			Description:
			> Adding docs to clone

			Tags:
			 * Development (%[2]s62ae28b72518aa18da2acb49%[2]s)

			$ %[1]s clone last -d "Adding examples to clone" --md
			ID: %[2]s62ae4b304ebb4f143c931d50%[2]s  
			Billable: %[2]syes%[2]s  
			Locked: %[2]sno%[2]s  
			Project: Clockify Cli (%[2]s621948458cb9606d934ebb1c%[2]s)  
			Task: Clone Command (%[2]s62ae4af04ebb4f143c931d2e%[2]s)  
			Interval: %[2]s2022-06-18 22:11:16%[2]s until %[2]snow%[2]s  
			Description:
			> Adding examples to clone

			Tags:
			 * Development (%[2]s62ae28b72518aa18da2acb49%[2]s)

			$ %[1]s clone last -d "Adding examples to in" -T pair --task "in command" --md
			ID: %[2]s62ae4dfe4ebb4f143c932106%[2]s
			Billable: %[2]syes%[2]s
			Locked: %[2]sno%[2]s
			Project: Clockify Cli (%[2]s621948458cb9606d934ebb1c%[2]s)
			Task: In Command (%[2]s62ae29e62518aa18da2acd14%[2]s)
			Interval: %[2]s2022-06-18 22:13:14%[2]s until %[2]snow%[2]s
			Description:
			> Adding examples to in

			Tags:
			 * Pair Programming (%[2]s621948708cb9606d934ebba7%[2]s)
		`, "clockify-cli", "`"),
		Args:      cmdutil.RequiredNamedArgs("time entry id"),
		ValidArgs: []string{timeentryhlp.AliasLast},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}
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

			id := strings.ToLower(strings.TrimSpace(args[0]))
			if id == timeentryhlp.AliasLast {
				id = timeentryhlp.AliasLatest
			}

			tec, err := timeentryhlp.GetTimeEntry(c, w, u, id)
			if err != nil {
				return err
			}

			tec.UserID = u
			tec.TimeInterval.End = nil

			noClosing, _ := cmd.Flags().GetBool("no-closing")

			dc := util.NewDescriptionCompleter(f)

			if tec, err = util.Do(
				tec,
				util.FillTimeEntryWithFlags(cmd.Flags()),
				func(tec dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
					if noClosing {
						return tec, nil
					}

					return util.ValidateClosingTimeEntry(f)(tec)
				},
				util.GetAllowNameForIDsFn(f.Config(), c),
				util.GetPropsInteractiveFn(c, dc, f.Config()),
				util.GetDatesInteractiveFn(f.Config()),
				util.GetValidateTimeEntryFn(f),
				func(tec dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
					if noClosing {
						return tec, nil
					}

					return util.OutInProgressFn(c)(tec)
				},
				util.CreateTimeEntryFn(c),
			); err != nil {
				return err
			}

			return util.PrintTimeEntryImpl(tec, f, cmd.OutOrStdout(), of)
		},
	}

	util.AddTimeEntryFlags(cmd, f, &of)
	util.AddTimeEntryDateFlags(cmd)
	cmd.Flags().BoolP("no-closing", "", false,
		"don't close any active time entry")

	return cmd
}
