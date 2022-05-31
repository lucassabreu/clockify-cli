package util

import (
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/client"
	"github.com/spf13/cobra"
)

// OutputFlags sets how to print out a list of clients
type OutputFlags struct {
	Format string
	CSV    bool
	JSON   bool
	Quiet  bool
}

func (of OutputFlags) Check() error {
	return cmdutil.XorFlag(map[string]bool{
		"format": of.Format != "",
		"json":   of.JSON,
		"csv":    of.CSV,
		"quiet":  of.Quiet,
	})
}

// AddReportFlags adds the default output flags for clients
func AddReportFlags(cmd *cobra.Command, of *OutputFlags) {
	cmd.Flags().StringVarP(&of.Format, "format", "f", "",
		"golang text/template format to be applied on each Client")
	cmd.Flags().BoolVarP(&of.JSON, "json", "j", false, "print as JSON")
	cmd.Flags().BoolVarP(&of.CSV, "csv", "v", false, "print as CSV")
	cmd.Flags().BoolVarP(&of.Quiet, "quiet", "q", false, "only display ids")
}

// Report prints out the clients
func Report(cs []dto.Client, out io.Writer, of OutputFlags) error {
	switch {
	case of.JSON:
		return output.ClientsJSONPrint(cs, out)
	case of.CSV:
		return output.ClientsCSVPrint(cs, out)
	case of.Format != "":
		return output.ClientPrintWithTemplate(of.Format)(cs, out)
	case of.Quiet:
		return output.ClientPrintQuietly(cs, out)
	default:
		return output.ClientPrint(cs, out)
	}
}
