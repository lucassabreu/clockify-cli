package defaults

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util/defaults"
	"gopkg.in/yaml.v3"
)

const (
	FORMAT_JSON = "json"
	FORMAT_YAML = "yaml"
)

type OutputFlags struct {
	Format string
}

// Report prints a DefaultTimeEntry using user's flags
func Report(of OutputFlags, out io.Writer, v defaults.DefaultTimeEntry) error {
	var b []byte
	switch of.Format {
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
