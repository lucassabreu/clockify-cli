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
	"os"

	"github.com/lucassabreu/clockify-cli/cmd/completion"
	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions for every session, execute once:

#### Linux:

` + "```" + `
$ clockify-cli completion bash > /etc/bash_completion.d/clockify-cli
` + "```" + `

#### MacOS:

` + "```" + `
$ clockify-cli completion bash > /usr/local/etc/bash_completion.d/clockify-cli
` + "```" + `

#### Zsh:

To load completions for each session, add this line to your ~/.zshrc:
` + "```" + `
source <(clockify-cli completion zsh)
` + "```" + `

You will need to start a new shell for this setup to take effect.

#### Fish:
To load completions for each session, execute once:
` + "```" + `
$ clockify-cli completion fish > ~/.config/fish/completions/clockify-cli.fish
` + "```",
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			return completion.GenZshCompletion(cmd, os.Stdout)
		case "fish":
			return cmd.Root().GenFishCompletion(os.Stdout, false)
		case "powershell":
			return cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
