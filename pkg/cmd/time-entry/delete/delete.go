package del

import (
	"errors"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timeentryhlp"
	"github.com/spf13/cobra"
)

// NewCmdDelete represents the delete command
func NewCmdDelete(f cmdutil.Factory) *cobra.Command {
	va := cmdcompl.ValidArgsSlide{timeentryhlp.AliasCurrent, timeentryhlp.AliasLast}
	cmd := &cobra.Command{
		Use: "delete { <time-entry-id> | " +
			va.IntoUseOptions() + " }...",
		Aliases:   []string{"del", "rm", "remove"},
		Args:      cmdutil.RequiredNamedArgs("time entry id"),
		ValidArgs: va.IntoValidArgs(),
		Short: `Delete time entry(ies), use id "` +
			timeentryhlp.AliasCurrent + `" to apply to time entry in progress`,
		Long: heredoc.Docf(`
			Delete time entries

			If you want to delete the current (running) time entry you can use "%s" instead of its ID.

			**Important**: this action can't be reverted, once the time entry is deleted its ID is lost.
		`,
			timeentryhlp.AliasCurrent,
		),
		Example: heredoc.Docf(`
			# trying to delete a time entry that does not exist, or from other workspace
			$ %[1]s 62af70d849445270d7c09fbc
			delete time entry "62af70d849445270d7c09fbc": TIMEENTRY with id 62af70d849445270d7c09fbc doesn't belong to WORKSPACE with id cccccccccccccccccccccccc (code: 501)

			# deleting the running time entry
			$ %[1]s current
			# no output

			# deleting the last time entry
			$ %[1]s last
			# no output

			# deleting multiple time entries
			$ %[1]s 62b5b51085815e619d7ae18d 62b5d55185815e619d7af928
			# no output

			# deleting last two entries
			$ %[1]s last last
			# no output
		`, "clockify-cli delete"),
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

			for i := range args {
				p := api.DeleteTimeEntryParam{
					Workspace:   w,
					TimeEntryID: args[i],
				}

				if p.TimeEntryID == timeentryhlp.AliasCurrent {
					te, err := c.GetTimeEntryInProgress(
						api.GetTimeEntryInProgressParam{
							Workspace: p.Workspace,
							UserID:    u,
						})

					if err != nil {
						return err
					}

					if te == nil {
						return errors.New("there is no time entry in progress")
					}

					p.TimeEntryID = te.ID
				}

				if p.TimeEntryID == timeentryhlp.AliasLast {
					te, err := timeentryhlp.GetLatestEntryEntry(c, p.Workspace, u)

					if err != nil {
						return err
					}

					p.TimeEntryID = te.ID
				}

				if err := c.DeleteTimeEntry(p); err != nil {
					return err
				}
			}

			return nil
		},
	}

	return cmd
}
