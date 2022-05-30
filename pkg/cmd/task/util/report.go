package util

import (
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/output/task"
	"github.com/spf13/cobra"
)

// OutputFlags
type OutputFlags struct {
	Format string
	JSON   bool
	CSV    bool
	Quiet  bool
}

// TaskAddReportFlags will add common format flags used for tasks
func TaskAddReportFlags(cmd *cobra.Command, of *OutputFlags) {
	cmd.Flags().StringVarP(&of.Format, "format", "f", "",
		"golang text/template format to be applied on each Client")
	cmd.Flags().BoolVarP(&of.JSON, "json", "j", false, "print as JSON")
	cmd.Flags().BoolVarP(&of.CSV, "csv", "v", false, "print as CSV")
	cmd.Flags().BoolVarP(&of.Quiet, "quiet", "q", false, "only display ids")
}

// TaskReport will output the task as set by the flags
func TaskReport(cmd *cobra.Command, of OutputFlags, tasks ...dto.Task) error {
	out := cmd.OutOrStdout()

	switch {
	case of.JSON:
		return task.TasksJSONPrint(tasks, out)
	case of.CSV:
		return task.TasksCSVPrint(tasks, out)
	case of.Quiet:
		return task.TaskPrintQuietly(tasks, out)
	case of.Format != "":
		return task.TaskPrintWithTemplate(of.Format)(tasks, out)
	default:
		return task.TaskPrint(tasks, out)
	}
}
