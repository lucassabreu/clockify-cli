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
	"fmt"

	"github.com/spf13/cobra"
)

var githubIssues bool
var trelloBoard bool

// projectSetCmd represents the projectSet command
var projectSetCmd = &cobra.Command{
	Use:   "set <project-name-or-id> <github-repo-or-trello-board>",
	Short: "Links a project with a GitHub:Issues' repository or Trello's Board",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("projectSet called")
	},
}

func init() {
	projectCmd.AddCommand(projectSetCmd)

	projectSetCmd.Flags().BoolVarP(&githubIssues, "github", "g", false, "link with GitHub:Issues")
	projectSetCmd.Flags().BoolVarP(&trelloBoard, "trello", "b", false, "link with Trello Board")
}
