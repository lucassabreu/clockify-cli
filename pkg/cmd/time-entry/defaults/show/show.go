package show

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util/defaults"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const formatYAML = "yaml"
const formatJSON = "json"

// NewCmdShow prints the default options for the current folder
func NewCmdShow(f cmdutil.Factory) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:  "show",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := f.TimeEntryDefaults().Read()
			if err != nil {
				return err
			}

			return report(cmd.OutOrStdout(), format, d)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", formatYAML, "output format")
	_ = cmdcompl.AddFixedSuggestionsToFlag(cmd, "format",
		cmdcompl.ValidArgsSlide{formatYAML, formatJSON})

	return cmd
}

func report(out io.Writer, format string, v defaults.DefaultTimeEntry) error {
	format = strings.ToLower(format)
	var b []byte
	switch format {
	case formatJSON:
		b, _ = json.Marshal(v)
	case formatYAML:
		b, _ = yaml.Marshal(v)
	default:
		return errors.New("invalid format")
	}

	_, err := out.Write(b)
	return err
}
