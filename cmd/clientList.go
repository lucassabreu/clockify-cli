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
)

// clientListCmd represents the list command
var clientListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List clients on Clockify",
	Aliases: []string{"ls"},
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		format, _ := cmd.Flags().GetString("format")
		asJSON, _ := cmd.Flags().GetBool("json")
		asCSV, _ := cmd.Flags().GetBool("csv")
		quiet, _ := cmd.Flags().GetBool("quiet")
		name, _ := cmd.Flags().GetString("name")

		workspace, err := getWorkspaceOrDefault(c)
		if err != nil {
			return err
		}

		p := api.GetClientsParam{
			Workspace:       workspace,
			Name:            name,
			PaginationParam: api.AllPages(),
		}

		if ok, _ := cmd.Flags().GetBool("not-archived"); ok {
			b := false
			p.Archived = &b
		}
		if ok, _ := cmd.Flags().GetBool("archived"); ok {
			b := true
			p.Archived = &b
		}

		clients, err := c.GetClients(p)
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

		return reportFn(clients, os.Stdout)
	}),
}

func init() {
	clientCmd.AddCommand(clientListCmd)

	clientListCmd.Flags().StringP("name", "n", "", "will be used to filter the tag by name")
	clientListCmd.Flags().BoolP("not-archived", "", false, "list only active projects")
	clientListCmd.Flags().BoolP("archived", "", false, "list only archived projects")
	clientListCmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on each Project")
	clientListCmd.Flags().BoolP("json", "j", false, "print as JSON")
	clientListCmd.Flags().BoolP("csv", "v", false, "print as CSV")
	clientListCmd.Flags().BoolP("quiet", "q", false, "only display ids")
}
