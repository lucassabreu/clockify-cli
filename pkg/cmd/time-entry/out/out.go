package out

import (
	"errors"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/time-entry"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
)

// NewCmdOut represents the out command
func NewCmdOut(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "out",
		Short: "Stops the last time entry",
		RunE: func(cmd *cobra.Command, args []string) error {
			var whenDate time.Time
			var err error

			whenString, _ := cmd.Flags().GetString("when")
			if whenDate, err = timehlp.ConvertToTime(whenString); err != nil {
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

			te, err := c.GetHydratedTimeEntryInProgress(
				api.GetTimeEntryInProgressParam{
					Workspace: w,
					UserID:    userID,
				})

			if te == nil && err == nil {
				return errors.New("no time entry in progress")
			}

			if err != nil {
				return err
			}

			if err = c.Out(api.OutParam{
				Workspace: w,
				UserID:    userID,
				End:       whenDate,
			}); err != nil {
				return err
			}

			te.TimeInterval.End = &whenDate

			return util.PrintTimeEntry(te,
				cmd, output.TimeFormatSimple, f.Config())
		},
	}

	util.AddPrintTimeEntriesFlags(cmd)

	cmd.Flags().String("when", time.Now().Format(timehlp.FullTimeFormat),
		"when the entry should be closed, "+
			"if not informed will use current time")

	return cmd
}
