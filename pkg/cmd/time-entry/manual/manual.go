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

package manual

import (
	"fmt"
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcomplutil"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/time-entry"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
)

// NewCmdManual represents the manual command
func NewCmdManual(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use: "manual [<project-id>] [<start>] [<end>] [<description>]",
		Short: "Creates a new completed time entry " +
			"(does not stop on-going time entries)",
		Args: cobra.MaximumNArgs(4),
		ValidArgsFunction: cmdcompl.CombineSuggestionsToArgs(
			cmdcomplutil.NewProjectAutoComplete(f)),
		RunE: func(cmd *cobra.Command, args []string) error {
			var whenToCloseDate time.Time
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
				tei.TimeInterval.Start, err = timehlp.ConvertToTime(args[1])
				if err != nil {
					return fmt.Errorf(
						"Fail to convert when to start: %w", err)
				}
			}

			if len(args) > 2 {
				whenToCloseDate, err = timehlp.ConvertToTime(args[2])
				if err != nil {
					return fmt.Errorf(
						"Fail to convert when to end: %w", err)
				}
				tei.TimeInterval.End = &whenToCloseDate
			}

			if len(args) > 3 {
				tei.Description = args[3]
			}

			dc := util.NewDescriptionCompleter(f)

			if tei, err = util.ManageEntry(
				tei,
				util.FillTimeEntryWithFlags(cmd.Flags()),
				func(tei dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
					if tei.TimeInterval.End != nil {
						return tei, nil
					}

					now, _ := timehlp.ConvertToTime(timehlp.NowTimeFormat)
					tei.TimeInterval.End = &now
					return tei, nil
				},
				util.GetAllowNameForIDsFn(f.Config(), c),
				util.GetPropsInteractiveFn(c, dc, f.Config()),
				util.GetDatesInteractiveFn(f.Config()),
				util.ValidateClosingTimeEntry(f),
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
