package split

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timeentryhlp"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewCmdSplit(
	f cmdutil.Factory,
	report func([]dto.TimeEntry, io.Writer, util.OutputFlags) error,
) *cobra.Command {
	of := util.OutputFlags{TimeFormat: timehlp.OnlyTimeFormat}
	va := cmdcompl.ValidArgsSlide{
		timeentryhlp.AliasCurrent,
		timeentryhlp.AliasLast,
		timeentryhlp.AliasLatest,
	}

	cmd := &cobra.Command{
		Use: "split { <time-entry-id> | " + va.IntoUseOptions() + " | ^n } " +
			" <time>...",
		Args: cobra.MatchAll(
			cmdutil.RequiredNamedArgs("time entry id"),
			cobra.MinimumNArgs(2),
		),
		ValidArgs: va.IntoValidArgs(),
		Short:     `Splits a time entry into multiple time entries`,
		Long: heredoc.Docf(`
			Split a time entry.
			The time arguments can be more than one, but must be increasing.

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
			$ %[1]s in --project cli --tag dev -d "Doing work before lunch" --task "edit" --md
			ID: %[2]s62ae4b304ebb4f143c931d50%[2]s  
			Billable: %[2]syes%[2]s  
			Locked: %[2]sno%[2]s  
			Project: Clockify Cli (%[2]s621948458cb9606d934ebb1c%[2]s)  
			Task: Edit Command (%[2]s62ae4af04ebb4f143c931d2e%[2]s)  
			Interval: %[2]s2022-06-18 11:01:16%[2]s until %[2]snow%[2]s  
			Description:
			> Adding docs to edit

			Tags:
			 * Development (%[2]s62ae28b72518aa18da2acb49%[2]s)

			# splits the time entry at lunch and now
			$ %[1]s split 12:00 13:30 --format '{{.ID}},{{.TimeInterval.Start|ft}}'
			62ae4b304ebb4f143c931d50,11:01
			3c931d502ae4b3064ebb4f14,12:00
			ebb4f143c962ae4b30431d50,13:30
		`, "clockify-cli", "`"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := of.Check(); err != nil {
				return err
			}

			splits := make([]time.Time, len(args)-1)
			for i := range splits {
				t, err := timehlp.ConvertToTime(args[1+i])
				if err != nil {
					return fmt.Errorf(
						"argument %d could not be converted to time: %w",
						i+2, err)
				}

				if i > 0 && t.Before(splits[i-1]) {
					return errors.New("splits must be in increasing order")
				}

				splits[i] = t
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

			te, err := timeentryhlp.GetTimeEntry(
				c,
				w,
				userID,
				args[0],
			)
			if err != nil {
				return err
			}

			if te.TimeInterval.Start.After(splits[0]) {
				return errors.New("time splits must be after " +
					te.TimeInterval.Start.Format(timehlp.FullTimeFormat))
			}

			if te.TimeInterval.End != nil && te.TimeInterval.End.Before(
				splits[len(splits)-1]) {
				return errors.New("time splits must be before " +
					te.TimeInterval.End.Format(timehlp.FullTimeFormat))
			}

			if _, err = c.UpdateTimeEntry(api.UpdateTimeEntryParam{
				Workspace:   te.WorkspaceID,
				TimeEntryID: te.ID,
				Description: te.Description,
				Start:       te.TimeInterval.Start,
				End:         &splits[0],
				Billable:    te.Billable,
				ProjectID:   te.ProjectID,
				TaskID:      te.TaskID,
				TagIDs:      te.TagIDs,
			}); err != nil {
				return err
			}

			tes := make([]dto.TimeEntry, len(splits)+1)
			getHydrated := func(i int, id string) error {
				t, err := c.GetHydratedTimeEntry(api.GetTimeEntryParam{
					TimeEntryID: id,
					Workspace:   w,
				})

				if err != nil {
					return err
				}
				tes[i] = *t
				return nil
			}

			eg := errgroup.Group{}
			eg.Go(func() error { return getHydrated(0, te.ID) })

			for i := range splits {
				i := i
				eg.Go(func() error {
					end := te.TimeInterval.End
					if i < len(splits)-1 {
						end = &splits[i+1]
					}

					te, err := c.CreateTimeEntry(api.CreateTimeEntryParam{
						Workspace:   te.WorkspaceID,
						Billable:    &te.Billable,
						Start:       splits[i],
						End:         end,
						ProjectID:   te.ProjectID,
						Description: te.Description,
						TagIDs:      te.TagIDs,
						TaskID:      te.TaskID,
					})

					if err != nil {
						return err
					}

					return getHydrated(i+1, te.ID)
				})

			}

			if err := eg.Wait(); err != nil {
				return err
			}

			return report(tes, cmd.OutOrStdout(), of)
		},
	}

	util.AddPrintTimeEntriesFlags(cmd, &of)
	util.AddPrintMultipleTimeEntriesFlags(cmd)

	return cmd
}
