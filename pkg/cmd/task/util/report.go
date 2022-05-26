package util

import (
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/output/task"
	"github.com/spf13/cobra"
)

// TaskAddReportFlags will add common format flags used for tasks
func TaskAddReportFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("format", "f", "",
		"golang text/template format to be applied on each Client")
	cmd.Flags().BoolP("json", "j", false, "print as JSON")
	cmd.Flags().BoolP("csv", "v", false, "print as CSV")
	cmd.Flags().BoolP("quiet", "q", false, "only display ids")
}

// TaskReport will output the task as set by the flags
func TaskReport(cmd *cobra.Command, tasks ...dto.Task) error {
	flag := cmd.Flags()
	out := cmd.OutOrStdout()

	if asJSON, _ := flag.GetBool("json"); asJSON {
		return task.TasksJSONPrint(tasks, out)
	}

	if asCSV, _ := flag.GetBool("csv"); asCSV {
		return task.TasksCSVPrint(tasks, out)
	}

	if format, _ := flag.GetString("format"); format != "" {
		return task.TaskPrintWithTemplate(format)(tasks, out)
	}

	if quiet, _ := flag.GetBool("quiet"); quiet {
		return task.TaskPrintQuietly(tasks, out)
	}

	return task.TaskPrint(tasks, out)
}
