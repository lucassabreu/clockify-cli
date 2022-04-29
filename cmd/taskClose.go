/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/spf13/cobra"
)

// taskCloseCmd represents the close command
var taskCloseCmd = &cobra.Command{
	Use:     "close <task>",
	Aliases: []string{"mark-as-done"},
	Short:   "Set the state of a task as DONE",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		c.AddClient()
	}),
}

func init() {
	taskCmd.AddCommand(taskCloseCmd)

	addProjectFlags(taskCloseCmd)
	taskAddReportFlags(taskCloseCmd)
}
