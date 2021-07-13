// Copyright © 2019 Lucas dos Santos Abreu <lucas.s.abreu@gmail.com>
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
	"github.com/spf13/viper"
)

// tagsCmd represents the tags command
var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List tags of workspace",
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		format, _ := cmd.Flags().GetString("format")
		quiet, _ := cmd.Flags().GetBool("quiet")
		archived, _ := cmd.Flags().GetBool("archived")
		name, _ := cmd.Flags().GetString("name")

		tags, err := getTags(c, name, archived)
		if err != nil {
			return err
		}

		var reportFn func([]dto.Tag, io.Writer) error

		reportFn = output.TagPrint
		if format != "" {
			reportFn = output.TagPrintWithTemplate(format)
		}

		if quiet {
			reportFn = output.TagPrintQuietly
		}

		return reportFn(tags, os.Stdout)
	}),
}

func getTags(c *api.Client, name string, archived bool) ([]dto.Tag, error) {
	return c.GetTags(api.GetTagsParam{
		Workspace:       viper.GetString(WORKSPACE),
		Name:            name,
		Archived:        archived,
		PaginationParam: api.PaginationParam{AllPages: true},
	})
}

func init() {
	rootCmd.AddCommand(tagsCmd)

	tagsCmd.Flags().StringP("name", "n", "", "will be used to filter the tag by name")
	tagsCmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on each Tag")
	tagsCmd.Flags().BoolP("quiet", "q", false, "only display ids")
	tagsCmd.Flags().BoolP("archived", "", false, "only display archived tags")
}
