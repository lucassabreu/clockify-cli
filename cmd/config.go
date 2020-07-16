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
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

const FORMAT_YAML = "yaml"
const FORMAT_JSON = "json"

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config [config-name]",
	Short: "Manages configuration file parameters",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")

		var b []byte

		var v interface{}
		if len(args) == 0 {
			v = viper.AllSettings()
		} else {
			v = viper.Get(args[0])
		}

		format = strings.ToLower(format)
		switch format {
		case FORMAT_JSON:
			b, _ = json.Marshal(v)

		case FORMAT_YAML:
			b, _ = yaml.Marshal(v)
		default:
			return errors.New("invalid format")
		}

		fmt.Println(string(b))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringP("format", "f", FORMAT_YAML, fmt.Sprintf("format of the output can be one of: %s, %s", FORMAT_YAML, FORMAT_JSON))
}
