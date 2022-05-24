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

package completion

import (
	"fmt"
	"io"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/spf13/cobra"
)

const (
	bash       = "bash"
	zsh        = "zsh"
	fish       = "fish"
	powershell = "powershell"
)

// completionCmd represents the completion command
func NewCmdCompletion() *cobra.Command {
	args := cmdcompl.ValidArgsSlide{bash, zsh, fish, powershell}

	cmd := &cobra.Command{
		Use:                   "completion " + args.IntoUse(),
		Short:                 "Generate completion script",
		DisableFlagsInUseLine: true,
		ValidArgs:             args.OnlyArgs(),
		Args:                  cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			switch strings.ToLower(args[0]) {
			case bash:
				return cmd.Root().GenBashCompletion(out)
			case zsh:
				return genZshCompletion(cmd, out)
			case fish:
				return cmd.Root().GenFishCompletion(out, false)
			case powershell:
				return cmd.Root().GenPowerShellCompletion(out)
			default:
				return nil
			}
		},
	}

	cmd.Long = heredoc.Docf(`
		To load completions for every session, execute once:

		#### Linux (Bash):

		%[1]s
		$ clockify-cli completion %[2]s > /etc/bash_cmdcompl.d/clockify-cli
		%[1]s

		#### Linux (Shell):

		%[1]s
		$ clockify-cli completion %[2]s > /etc/bash_cmdcompl.d/clockify-cli
		%[1]s

		#### MacOS:

		%[1]s
		$ clockify-cli completion %[2]s > /usr/local/etc/bash_cmdcompl.d/clockify-cli
		%[1]s

		#### Zsh:

		To load completions for each session, add this line to your ~/.zshrc:
		%[1]s
		source <(clockify-cli completion %[3]s)
		%[1]s

		You will need to start a new shell for this setup to take effect.

		#### Fish:
		To load completions for each session, execute once:
		%[1]s
		$ clockify-cli completion %[4]s > ~/.config/fish/completions/clockify-cli.fish
		%[1]s`, "```", bash, zsh, fish)

	return cmd
}

func genZshCompletion(cmd *cobra.Command, w io.Writer) error {
	if _, err := fmt.Fprintln(w,
		"autoload -U compinit; compinit"); err != nil {
		return err
	}

	if err := cmd.Root().GenZshCompletion(w); err != nil {
		return err
	}

	_, err := fmt.Fprintln(w, "compdef _clockify-cli clockify-cli")
	return err
}
