package invoiced

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/time-entry"
	"github.com/lucassabreu/clockify-cli/pkg/timeentryhlp"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/spf13/cobra"
)

// NewCmdInvoiced represents invoiced command
func NewCmdInvoiced(f cmdutil.Factory) []*cobra.Command {
	of := util.OutputFlags{TimeFormat: output.TimeFormatSimple}
	addCmd := func(cmd *cobra.Command) *cobra.Command {
		util.AddPrintTimeEntriesFlags(cmd, &of)
		util.AddPrintMultipleTimeEntriesFlags(cmd)

		return cmd
	}

	va := cmdcompl.ValidArgsSlide{
		timeentryhlp.AliasLast, timeentryhlp.AliasCurrent}
	use := "{ <time-entry-id> | " + va.IntoUseOptions() + " }..."
	args := cmdutil.RequiredNamedArgs("time entry id")

	return []*cobra.Command{
		addCmd(&cobra.Command{
			Use:   "mark-invoiced " + use,
			Short: "Marks times entries as invoiced",
			Long: "Marks times entries as invoiced\n\n" +
				util.HelpMoreInfoAboutPrinting,
			Example: heredoc.Docf(`
				# when the workspace does not allow invoicing
				$ %[1]s 62b49641f4b27f4ed7d20e75
				Forbidden (code: 403)

				# set the running time entry as invoiced
				$ %[1]s current --quiet
				62b49641f4b27f4ed7d20e75

				# setting multiple time entries as invoiced
				$ %[1]s 62b5b51085815e619d7ae18d 62b5d55185815e619d7af928 --quiet
				62b5b51085815e619d7ae18d
				62b5d55185815e619d7af928
			`, "clockify-cli mark-invoiced"),
			Args:      args,
			ValidArgs: va,
			RunE:      changeInvoiced(f, &of, true),
		}),
		addCmd(&cobra.Command{
			Use:   "mark-not-invoiced " + use,
			Short: "Mark times entries as not invoiced",
			Long: "Mark times entries as not invoiced\n\n" +
				util.HelpMoreInfoAboutPrinting,
			Example: heredoc.Docf(`
				# when the workspace does not allow invoicing
				$ %[1]s 62b49641f4b27f4ed7d20e75
				Forbidden (code: 403)

				# set the running time entry as not invoiced
				$ %[1]s current --quiet
				62b49641f4b27f4ed7d20e75

				# setting multiple time entries as not invoiced
				$ %[1]s 62b5b51085815e619d7ae18d 62b5d55185815e619d7af928 --quiet
				62b5b51085815e619d7ae18d
				62b5d55185815e619d7af928
			`, "clockify-cli mark-not-invoiced"),
			Args:      args,
			ValidArgs: va,
			RunE:      changeInvoiced(f, &of, false),
		}),
	}
}

func changeInvoiced(
	f cmdutil.Factory, of *util.OutputFlags, invoiced bool,
) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
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

		args = strhlp.Unique(args)
		tes := make([]dto.TimeEntry, len(args))
		for i, id := range args {
			if id == timeentryhlp.AliasCurrent ||
				id == timeentryhlp.AliasLast {
				tei, err := timeentryhlp.GetTimeEntry(c, w, u, id)
				if err != nil {
					return err
				}
				id = tei.ID
				args[i] = id
			}

			te, err := c.GetHydratedTimeEntry(api.GetTimeEntryParam{
				Workspace:   w,
				TimeEntryID: id,
			})
			if err != nil {
				return err
			}

			tes[i] = *te
		}

		if err := c.ChangeInvoiced(api.ChangeInvoicedParam{
			Workspace:    w,
			TimeEntryIDs: args,
			Invoiced:     invoiced,
		}); err != nil {
			return err
		}

		return util.PrintTimeEntries(tes, cmd.OutOrStdout(), f.Config(), *of)
	}
}
