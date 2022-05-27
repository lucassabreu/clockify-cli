package util

import (
	"io"
	"os"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/output/project"
	"github.com/spf13/cobra"
)

// OutputFlags defines how to print the project
type OutputFlags struct {
	JSON   bool
	CSV    bool
	Quiet  bool
	Format string
}

func (of OutputFlags) Check() error {
	return cmdutil.XorFlag(map[string]bool{
		"format": of.Format != "",
		"json":   of.JSON,
		"csv":    of.CSV,
		"quiet":  of.Quiet,
	})
}

// AddReportFlags adds the common flags to print projects
func AddReportFlags(cmd *cobra.Command, of *OutputFlags) {
	cmd.Flags().StringVarP(&of.Format, "format", "f", "",
		"golang text/template format to be applied on each Project")
	cmd.Flags().BoolVarP(&of.JSON, "json", "j", false, "print as JSON")
	cmd.Flags().BoolVarP(&of.CSV, "csv", "v", false, "print as CSV")
	cmd.Flags().BoolVarP(&of.Quiet, "quiet", "q", false, "only display ids")
}

// Report will print the projects as set by the flags
func Report(list []dto.Project, out io.Writer, f OutputFlags) error {
	switch {
	case f.JSON:
		return project.ProjectsJSONPrint(list, out)
	case f.CSV:
		return project.ProjectsCSVPrint(list, out)
	case f.Quiet:
		return project.ProjectPrintQuietly(list, out)
	case f.Format != "":
		return project.ProjectPrintWithTemplate(f.Format)(list, out)
	default:
		return project.ProjectPrint(list, os.Stdout)
	}
}
