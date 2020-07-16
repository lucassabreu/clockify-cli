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
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "clockify-cli",
	Short:         "Allow to integrate with Clockify through terminal",
	SilenceErrors: true,
	SilenceUsage:  true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if viper.GetBool("debug") {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
		} else {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.clockify-cli.yaml)")

	rootCmd.PersistentFlags().StringP("token", "t", "", `clockify's token (defaults to env $CLOCKIFY_TOKEN)
	Can be generated here: https://clockify.me/user/settings#generateApiKeyBtn`)
	_ = viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))

	rootCmd.PersistentFlags().StringP("workspace", "w", "", "workspace to be used (defaults to env $CLOCKIFY_WORKSPACE)")
	_ = viper.BindPFlag("workspace", rootCmd.PersistentFlags().Lookup("workspace"))

	rootCmd.PersistentFlags().StringP("user-id", "u", "", "user id from the token (defaults to env $CLOCKIFY_USER_ID)")
	_ = viper.BindPFlag("user.id", rootCmd.PersistentFlags().Lookup("user-id"))

	rootCmd.PersistentFlags().Bool("debug", false, "show debug log (defaults to env $CLOCKIFY_DEBUG)")
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.PersistentFlags().BoolP("interactive", "i", false, "show interactive log (defaults to env $CLOCKIFY_INTERACTIVE)")
	_ = viper.BindPFlag("interactive", rootCmd.PersistentFlags().Lookup("interactive"))

	_ = rootCmd.MarkFlagRequired("token")

	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println(cmd.UsageString())
		return nil
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".clockify-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".clockify-cli")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetEnvPrefix("clockify")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			return
		default:
			fmt.Print(err)
			return
		}
	}
}
