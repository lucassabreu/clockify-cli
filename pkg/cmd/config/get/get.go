package get

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const FORMAT_YAML = "yaml"
const FORMAT_JSON = "json"

func NewCmdGet(
	f cmdutil.Factory,
	validParameters cmdcompl.ValidArgsMap,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [param]",
		Short: "Retrieves one or all parameters set by the user",
		Long: heredoc.Doc(`
			Retrieves one or all parameters set by the user.

			Available parameters are:
		`) + validParameters.Long(),
		Args:      cobra.MaximumNArgs(1),
		ValidArgs: validParameters.IntoValidArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			format, err := cmd.Flags().GetString("format")
			if err != nil {
				return err
			}

			param := ""
			if len(args) > 0 {
				param = args[0]
			}

			v := run(f.Config(), param)
			return show(cmd.OutOrStdout(), format, v)
		},
	}

	cmd.Flags().StringP("format", "f", FORMAT_YAML,
		"output format (when not setting or initializing)")
	_ = cmdcompl.AddFixedSuggestionsToFlag(cmd, "format",
		cmdcompl.ValidArgsSlide{FORMAT_YAML, FORMAT_JSON})

	return cmd
}

func run(config cmdutil.Config, param string) interface{} {
	if param == "" {
		return config.All()
	}

	return config.Get(param)
}

func show(out io.Writer, format string, v interface{}) error {
	format = strings.ToLower(format)
	var b []byte
	switch format {
	case FORMAT_JSON:
		b, _ = json.Marshal(v)
	case FORMAT_YAML:
		b, _ = yaml.Marshal(v)
	default:
		return errors.New("invalid format")
	}

	_, err := out.Write(b)
	return err
}
