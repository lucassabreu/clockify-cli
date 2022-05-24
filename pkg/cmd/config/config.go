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

package config

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/config/get"
	initialize "github.com/lucassabreu/clockify-cli/pkg/cmd/config/init"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/config/set"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"

	"github.com/spf13/cobra"
)

var validParameters = cmdcompl.ValidArgsMap{
	cmdutil.CONF_TOKEN:     "clockify's token",
	cmdutil.CONF_WORKSPACE: "workspace to be used",
	cmdutil.CONF_USER_ID:   "user id from the token",
	cmdutil.CONF_ALLOW_NAME_FOR_ID: "allow to input the name of the entity " +
		"instead of its ID (projects, clients, tasks, users and tags)",
	cmdutil.CONF_INTERACTIVE: "show interactive mode",
	cmdutil.CONF_WORKWEEK_DAYS: "days of the week were your expected to " +
		"work (use comma to set multiple)",
	cmdutil.CONF_ALLOW_INCOMPLETE: "should allow starting time entries with " +
		"missing required values",
	cmdutil.CONF_SHOW_TASKS: "should show an extra column with the task " +
		"description",
	cmdutil.CONF_DESCR_AUTOCOMP: "autocomplete description looking at " +
		"recent time entries",
	cmdutil.CONF_DESCR_AUTOCOMP_DAYS: "how many days should be considered " +
		"for the description autocomplete",
	cmdutil.CONF_SHOW_TOTAL_DURATION: "adds a totals line on time entry " +
		"reports with the sum of the time entries duration",
	cmdutil.CONF_DEBUG: "if extra logs should be shown",
}

// NewCmdConfig represents the config command
func NewCmdConfig(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <command>",
		Short: "Manages configuration file parameters",
		Args:  cobra.MaximumNArgs(0),
	}

	cmd.AddCommand(initialize.NewCmdInit(f))
	cmd.AddCommand(set.NewCmdSet(f, validParameters))
	cmd.AddCommand(get.NewCmdGet(f, validParameters))

	return cmd
}
