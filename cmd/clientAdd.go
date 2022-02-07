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
	"os"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// clientAddCmd represents the add command
var clientAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"new", "create"},
	Short:   "Adds a client to the Clockify workspace",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		format, _ := cmd.Flags().GetString("format")
		asJSON, _ := cmd.Flags().GetBool("json")
		asCSV, _ := cmd.Flags().GetBool("csv")
		quiet, _ := cmd.Flags().GetBool("quiet")
		name, _ := cmd.Flags().GetString("name")

		workspace := viper.GetString(WORKSPACE)
		client, err := c.AddClient(api.AddClientParam{
			Workspace: workspace,
			Name:      name,
		})
		if err != nil {
			return err
		}

		var reportFn func([]dto.Client, io.Writer) error

		reportFn = output.ClientPrint
		if asJSON {
			reportFn = output.ClientsJSONPrint
		}

		if asCSV {
			reportFn = output.ClientsCSVPrint
		}

		if format != "" {
			reportFn = output.ClientPrintWithTemplate(format)
		}

		if quiet {
			reportFn = output.ClientPrintQuietly
		}

		return reportFn([]dto.Client{client}, os.Stdout)
	}),
}

func init() {
	clientCmd.AddCommand(clientAddCmd)

	clientAddCmd.Flags().StringP("name", "n", "", "the name of the new client")
	clientAddCmd.MarkFlagRequired("name")

	clientAddCmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on each Client")
	clientAddCmd.Flags().BoolP("json", "j", false, "print as JSON")
	clientAddCmd.Flags().BoolP("csv", "v", false, "print as CSV")
	clientAddCmd.Flags().BoolP("quiet", "q", false, "only display ids")
}
