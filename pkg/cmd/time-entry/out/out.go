package out

import (
	"errors"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/time-entry"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
)

// NewCmdOut represents the out command
func NewCmdOut(f cmdutil.Factory) *cobra.Command {
	of := util.OutputFlags{TimeFormat: output.TimeFormatSimple}
	cmd := &cobra.Command{
		Use:   "out",
		Short: "Stops the running time entry",
		Long: heredoc.Docf(`
			Stops the running time entry.

			If no value is set on %[1]s--when%[1]s, then current time will be used.

			When setting the end time you can use any of the following formats to set it:
			%[2]s

			Use %[1]sclockify-cli edit current%[1]s to edit any properties before ending it.
			%[3]s
		`, "`",
			util.HelpDateTimeFormats,
			util.HelpMoreInfoAboutPrinting,
		),
		Example: heredoc.Docf(`
			# stop running time entry with current time
			$ %[1]s out --md
			ID: %[2]s62af6b0f4ebb4f143c94880e%[2]s  
			Billable: %[2]syes%[2]s  
			Locked: %[2]sno%[2]s  
			Project: Clockify Cli (%[2]s621948458cb9606d934ebb1c%[2]s)  
			Task: Out Command (%[2]s62af66454ebb4f143c948263%[2]s)  
			Interval: %[2]s2022-06-19 18:29:32%[2]s until %[2]s2022-06-19 18:52:13%[2]s  
			Description:
			> Adding examples

			Tags:
			 * Development (%[2]s62ae28b72518aa18da2acb49%[2]s)

			# clone last and stopping it in 10 minutes
			$ %[1]s clone last -i=0 -d 'More examples' -q
			62af70d849445270d7c09fbd

			$ %[1]s out --when +10m --md
			ID: %[2]s62af70d849445270d7c09fbd%[2]s  
			Billable: %[2]syes%[2]s  
			Locked: %[2]sno%[2]s  
			Project: Clockify Cli (%[2]s621948458cb9606d934ebb1c%[2]s)  
			Task: Out Command (%[2]s62af666349445270d7c09285%[2]s)  
			Interval: %[2]s2022-06-19 18:54:12%[2]s until %[2]s2022-06-19 19:08:26%[2]s  
			Description:
			> More examples

			Tags:
			 * Development (%[2]s62ae28b72518aa18da2acb49%[2]s)
		`, "clockify-cli", "`"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

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

			return util.PrintTimeEntry(te, cmd.OutOrStdout(), f.Config(), of)
		},
	}

	util.AddPrintTimeEntriesFlags(cmd, &of)

	cmd.Flags().String("when", time.Now().Format(timehlp.FullTimeFormat),
		"when the entry should be closed, "+
			"if not informed will use current time")

	return cmd
}
