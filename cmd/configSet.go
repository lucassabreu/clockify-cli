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

	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configSetCmd represents the configSet command
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Updates or adds a key to the file",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		var value *string
		if len(args) >= 2 {
			value = &args[1]
		}

		if !viper.GetBool("interactive") && value == nil {
			return fmt.Errorf("you must inform the value of the config %s", key)
		}

		if value == nil {
			v := ""
			err := survey.AskOne(
				&survey.Input{
					Message: fmt.Sprintf("Value for %s:", key),
					Default: viper.GetString(key),
				},
				&v,
				nil,
			)

			if err != nil {
				return err
			}

			value = &v
		}

		viper.Set(key, value)
		return saveConfigFile()
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
}
