package util

import (
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/output/user"
	"github.com/spf13/cobra"
)

// OutputFlags sets how to print out a list of users
type OutputFlags struct {
	Format string
	JSON   bool
	Quiet  bool
}

func (of OutputFlags) Check() error {
	return cmdutil.XorFlag(map[string]bool{
		"format": of.Format != "",
		"json":   of.JSON,
		"quiet":  of.Quiet,
	})
}

// AddReportFlags adds the default output flags for users
func AddReportFlags(cmd *cobra.Command, of *OutputFlags) {
	cmd.Flags().StringVarP(&of.Format, "format", "f", "",
		"golang text/template format to be applied on each workspace")
	cmd.Flags().BoolVarP(&of.Quiet, "quiet", "q", false, "only display ids")
	cmd.Flags().BoolVarP(&of.JSON, "json", "j", false, "print as json")
}

// Report prints out the users
func Report(u []dto.User, out io.Writer, of OutputFlags) error {
	switch {
	case of.Format != "":
		return user.UserPrintWithTemplate(of.Format)(u, out)
	case of.Quiet:
		return user.UserPrintQuietly(u, out)
	default:
		return user.UserPrint(u, out)
	}
}
