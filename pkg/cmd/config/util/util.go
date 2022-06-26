package util

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const FormatYAML = "yaml"
const FormatJSON = "json"

// AddReportFlags adds the format flag
func AddReportFlags(cmd *cobra.Command, format *string) error {
	cmd.Flags().StringVarP(format, "format", "f", FormatYAML, "output format")
	return cmdcompl.AddFixedSuggestionsToFlag(cmd, "format",
		cmdcompl.ValidArgsSlide{FormatYAML, FormatJSON})
}

// Report prints the value as YAML or JSON
func Report(out io.Writer, format string, v interface{}) error {
	format = strings.ToLower(format)
	var b []byte
	switch format {
	case FormatJSON:
		b, _ = json.Marshal(v)
	case FormatYAML:
		b, _ = yaml.Marshal(v)
	default:
		return errors.New("invalid format")
	}

	_, err := out.Write(b)
	return err
}
