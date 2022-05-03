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
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/output"
	"github.com/spf13/cobra"
)

// taskCmd represents the client command
var taskCmd = &cobra.Command{
	Use:     "task",
	Aliases: []string{"tasks"},
	Short:   "List/add tasks of/to a project",
}

func init() {
	rootCmd.AddCommand(taskCmd)
}

func taskReport(cmd *cobra.Command, tasks ...dto.Task) error {
	var reportFn func([]dto.Task, io.Writer) error

	reportFn = output.TaskPrint
	if asJSON, _ := cmd.Flags().GetBool("json"); asJSON {
		reportFn = output.TasksJSONPrint
	}

	if asCSV, _ := cmd.Flags().GetBool("csv"); asCSV {
		reportFn = output.TasksCSVPrint
	}

	if format, _ := cmd.Flags().GetString("format"); format != "" {
		reportFn = output.TaskPrintWithTemplate(format)
	}

	if quiet, _ := cmd.Flags().GetBool("quiet"); quiet {
		reportFn = output.TaskPrintQuietly
	}

	return reportFn(tasks, cmd.OutOrStdout())
}

func taskAddReportFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on each Client")
	cmd.Flags().BoolP("json", "j", false, "print as JSON")
	cmd.Flags().BoolP("csv", "v", false, "print as CSV")
	cmd.Flags().BoolP("quiet", "q", false, "only display ids")
}
