package set

import (
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/spf13/cobra"
)

// NewCmdSet will update the value of one parameter
func NewCmdSet(
	f cmdutil.Factory,
	validParameters cmdcompl.ValidArgsMap,
) *cobra.Command {
	cmd := &cobra.Command{
		Use: "set <param> <value>",
		Args: cobra.MatchAll(
			cmdutil.RequiredNamedArgs("param", "value"),
			cobra.ExactArgs(2),
		),
		ValidArgs: validParameters.IntoValidArgs(),
		Short:     "Changes the value of one parameter",
		Example: heredoc.Docf(`
			$ %[1]s token "Yamdas569"
			$ %[1]s workweek-days monday,tuesday,wednesday,thursday,friday
			$ %[1]s show-task true
			$ %[1]s user.id 4564d5a6s4d54a5s4dasd5
		`, "clockify-cli config set"),
		RunE: func(cmd *cobra.Command, args []string) error {
			param := args[0]
			value := args[1]
			config := f.Config()

			switch param {
			case cmdutil.CONF_WORKWEEK_DAYS:
				ws := strings.Split(strings.ToLower(value), ",")
				ws = strhlp.Filter(
					func(s string) bool {
						return strhlp.Search(s, cmdutil.GetWeekdays()) != -1
					},
					ws,
				)
				config.SetStringSlice(param, ws)
			default:
				config.SetString(param, value)
			}

			return config.Save()
		},
	}

	return cmd
}
