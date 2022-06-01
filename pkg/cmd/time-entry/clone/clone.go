package clone

import (
	"strings"

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
		Use: "clone [" + timeentryhlp.AliasLast + "|<time-entry-id>]",
		Short: "Copy a time entry and starts it " +
			`(use "` + timeentryhlp.AliasLast + `" to copy the last one)`,
		Args:      cobra.ExactArgs(1),
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
