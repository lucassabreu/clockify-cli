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
	"strings"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/reports"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// tagsCmd represents the tags command
var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List tags of workspace",
	Run: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) {
		format, _ := cmd.Flags().GetString("format")
		quiet, _ := cmd.Flags().GetBool("quiet")

		tags, err := c.GetTags(api.GetTagsParam{
			Workspace: viper.GetString("workspace"),
		})

		if err != nil {
			printError(err)
			return
		}

		name, _ := cmd.Flags().GetString("name")
		tags = filterTags(name, tags)

		var reportFn func([]dto.Tag, io.Writer) error

		reportFn = reports.TagPrint
		if format != "" {
			reportFn = reports.TagPrintWithTemplate(format)
		}

		if quiet {
			reportFn = reports.TagPrintQuietly
		}

		if err = reportFn(tags, os.Stdout); err != nil {
			printError(err)
		}
	}),
}

func filterTags(name string, tags []dto.Tag) []dto.Tag {
	if name == "" {
		return tags
	}

	ts := make([]dto.Tag, 0)

	for _, t := range tags {
		if strings.Contains(strings.ToLower(t.Name), strings.ToLower(name)) {
			ts = append(ts, t)
		}
	}

	return ts
}

func init() {
	rootCmd.AddCommand(tagsCmd)

	tagsCmd.Flags().StringP("name", "n", "", "will be used to filter the tag by name")
	tagsCmd.Flags().StringP("format", "f", "", "golang text/template format to be applyed on each Tag")
	tagsCmd.Flags().BoolP("quiet", "q", false, "only display ids")
}
