package me

import (
	"io"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	output "github.com/lucassabreu/clockify-cli/pkg/output/user"
	"github.com/spf13/cobra"
)

// NewCmdMe represents the me command
func NewCmdMe(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "me",
		Short: "Show the user's token info",
		RunE: func(cmd *cobra.Command, args []string) error {
			format, err := cmd.Flags().GetString("format")
			if err != nil {
				return err
			}
			asJSON, err := cmd.Flags().GetBool("json")
			if err != nil {
				return err
			}
			quiet, err := cmd.Flags().GetBool("quiet")
			if err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			u, err := run(c)
			if err != nil {
				return err
			}

			return report(u, cmd.OutOrStdout(), format, asJSON, quiet)
		},
	}

	cmd.Flags().StringP("format", "f", "",
		"golang text/template format to be applied on the user")
	cmd.Flags().BoolP("json", "j", false, "print as json")
	cmd.Flags().BoolP("quiet", "q", false, "only display ids")

	return cmd
}

func run(c *api.Client) (dto.User, error) {
	return c.GetMe()
}

func report(u dto.User, out io.Writer,
	format string, asJSON, quiet bool) error {

	if asJSON {
		return output.UserJSONPrint(u, out)
	}

	list := []dto.User{u}
	if format != "" {
		return output.UserPrintWithTemplate(format)(list, out)
	}

	if quiet {
		return output.UserPrintQuietly(list, out)
	}

	return output.UserPrint(list, out)
}
