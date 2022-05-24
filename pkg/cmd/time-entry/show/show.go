/*
Copyright © 2021 Lucas dos Santos Abreu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package show

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timeentryhlp"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
)

// NewCmdShow represents the show command
func NewCmdShow(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use: "show [" + timeentryhlp.AliasCurrent + "|" +
			timeentryhlp.AliasLast + "|<time-entry-id>|^n]",
		ValidArgs: []string{timeentryhlp.AliasCurrent, timeentryhlp.AliasLast},
		Args:      cobra.MaximumNArgs(1),
		Short:     "Show detailed information about one time entry.",
		Long: heredoc.Doc(`
			Show detailed information about one time entry.
			Shows current one by default
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			userID, err := f.GetUserID()
			if err != nil {
				return err
			}

			w, err := f.GetWorkspaceID()
			if err != nil {
				return err
			}

			id := timeentryhlp.AliasCurrent
			if len(args) > 0 {
				id = args[0]
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			tei, err := timeentryhlp.GetTimeEntry(c, w, userID, id)
			if err != nil {
				return err
			}

			return util.PrintTimeEntryImpl(tei,
				f, cmd, timehlp.FullTimeFormat)
		},
	}

	util.AddPrintTimeEntriesFlags(cmd)
	_ = cmd.MarkFlagRequired("workspace")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}
