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

package in

import (
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcomplutil"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/time-entry"

	"github.com/spf13/cobra"
)

// NewCmdIn represents the in command
func NewCmdIn(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use: "in [<project-id>] [<description>]",
		Short: "Create a new time entry and starts it " +
			"(will close time entries not closed)",
		Args: cobra.MaximumNArgs(2),
		ValidArgsFunction: cmdcompl.CombineSuggestionsToArgs(
			cmdcomplutil.NewProjectAutoComplete(f)),
		Aliases: []string{"start"},
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			tei := dto.TimeEntryImpl{}

			if tei.WorkspaceID, err = f.GetWorkspaceID(); err != nil {
				return err
			}

			if tei.UserID, err = f.GetUserID(); err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			if len(args) > 0 {
				tei.ProjectID = args[0]
			}

			if len(args) > 1 {
				tei.Description = args[1]
			}

			dc := util.NewDescriptionCompleter(f)

			if tei, err = util.ManageEntry(
				tei,
				util.FillTimeEntryWithFlags(cmd.Flags()),
				util.ValidateClosingTimeEntry(f),
				util.GetAllowNameForIDsFn(f.Config(), c),
				util.GetPropsInteractiveFn(c, dc, f.Config()),
				util.GetDatesInteractiveFn(f.Config()),
				util.GetValidateTimeEntryFn(f),
				util.OutInProgressFn(c),
				util.CreateTimeEntryFn(c),
			); err != nil {
				return err
			}

			return util.PrintTimeEntryImpl(tei,
				f, cmd, output.TimeFormatSimple)
		},
	}

	util.AddTimeEntryFlags(cmd, f)
	util.AddTimeEntryDateFlags(cmd)

	return cmd
}