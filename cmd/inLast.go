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

	"github.com/spf13/cobra"
)

// inLastCmd represents the last command
var inLastCmd = &cobra.Command{
	Use:   "last",
	Short: "Copy the last time entry and starts it",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("in last called")
	},
}

func init() {
	inCmd.AddCommand(inLastCmd)
	inLastCmd.Flags().String("when", "", "when the entry should be closed, if not informed will use current time")
}
