// Copyright © 2020 Lucas dos Santos Abreu <lucas.s.abreu@gmail.com>
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

package del

import (
	"errors"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timeentryhlp"
	"github.com/spf13/cobra"
)

// NewCmdDelete represents the delete command
func NewCmdDelete(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete [" + timeentryhlp.AliasCurrent +
			"|<time-entry-id>]...",
		Aliases:   []string{"del", "rm", "remove"},
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: []string{timeentryhlp.AliasCurrent},
		Short: `Delete time entry(ies), use id "` +
			timeentryhlp.AliasCurrent + `" to apply to time entry in progress`,
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

				if err := c.DeleteTimeEntry(p); err != nil {
					return err
				}
			}

			return nil
		},
	}

	return cmd
}
