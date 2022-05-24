// Copyright © 2019 Lucas dos Santos Abreu <lucas.s.abreu@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package invoiced

import (
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/time-entry"
	"github.com/lucassabreu/clockify-cli/pkg/timeentryhlp"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/spf13/cobra"
)

// NewCmdInvoiced represents invoiced command
func NewCmdInvoiced(f cmdutil.Factory) (cmds []*cobra.Command) {

	addCmd := func(cmd *cobra.Command) {
		util.AddPrintTimeEntriesFlags(cmd)
		util.AddPrintMultipleTimeEntriesFlags(cmd)

		cmds = append(cmds, cmd)
	}

	use := "[" + timeentryhlp.AliasCurrent + "|" + timeentryhlp.AliasLast +
		"|<time-entry-id>]..."

	va := []string{timeentryhlp.AliasLast, timeentryhlp.AliasCurrent}

	addCmd(&cobra.Command{
		Use:       "mark-invoiced " + use,
		Short:     "Marks times entries as invoiced",
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: va,
		RunE:      changeInvoiced(f, true),
	})

	addCmd(&cobra.Command{
		Use:       "mark-not-invoiced " + use,
		Short:     "Mark times entries as not invoiced",
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: va,
		RunE:      changeInvoiced(f, false),
	})

	return cmds
}

func changeInvoiced(
	f cmdutil.Factory, invoiced bool,
) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
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

		return util.PrintTimeEntries(tes,
			cmd, output.TimeFormatSimple, f.Config())
	}
}
