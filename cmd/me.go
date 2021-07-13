// Copyright © 2020 Lucas dos Santos Abreu <lucas.s.abreu@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"io"
	"os"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/output"
	"github.com/spf13/cobra"
)

// meCmd represents the me command
var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Show the user info",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		format, _ := cmd.Flags().GetString("format")
		asJSON, _ := cmd.Flags().GetBool("json")

		u, err := c.GetMe()
		if err != nil {
			return err
		}

		var reportFn func(dto.User, io.Writer) error
		reportFn = func(u dto.User, w io.Writer) error {
			return output.UserPrint([]dto.User{u}, w)
		}

		if asJSON {
			reportFn = output.UserJSONPrint
		}

		if format != "" {
			reportFn = func(u dto.User, w io.Writer) error {
				return output.UserPrintWithTemplate(format)(
					[]dto.User{u},
					w,
				)
			}
		}

		return reportFn(u, os.Stdout)
	}),
}

func init() {
	rootCmd.AddCommand(meCmd)

	meCmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on the user")
	meCmd.Flags().BoolP("json", "j", false, "print as json")
}
