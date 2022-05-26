package set

import (
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/spf13/cobra"
)

func NewCmdSet(
	f cmdutil.Factory,
	validParameters cmdcompl.ValidArgsMap,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "set [param] [value]",
		Args:      cobra.ExactArgs(2),
		ValidArgs: validParameters.IntoValidArgs(),
		Short:     "Change the value of one parameter",
		Long: heredoc.Doc(`
			Change the value of one parameter.
			Available parameters are:
		`) + validParameters.Long(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(f.Config(), args[0], args[1])
		},
	}

	return cmd
}
func run(config cmdutil.Config, param, value string) error {
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
}
