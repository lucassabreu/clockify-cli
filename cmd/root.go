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
	"log"
	"os"

	"github.com/lucassabreu/clockify-cli/api"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var token string
var workspace string
var githubToken string
var trelloToken string
var debug bool
var userID string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "clockify-cli",
	Short: "Allow to integrate with Clockify through terminal",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.clockify-cli.yaml)")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "clockify's token, can be generated here: https://clockify.me/user/settings#generateApiKeyBtn")
	rootCmd.PersistentFlags().StringVarP(&workspace, "workspace", "w", "", "workspace to be used")
	rootCmd.PersistentFlags().StringVarP(&userID, "user-id", "u", "", "user id from the token")

	rootCmd.PersistentFlags().StringVar(&githubToken, "github-token", "", "gitHub's token")
	rootCmd.PersistentFlags().StringVar(&trelloToken, "trello-token", "", "trello's token")

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "show debug log")

	rootCmd.MarkFlagRequired("token")
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

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func withClockifyClient(fn func(cmd *cobra.Command, args []string, c *api.Client)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		c, err := api.NewClient(token)
		if err != nil {
			printError(err)
			return
		}

		if debug {
			c.SetDebugLogger(
				log.New(os.Stdout, "DEBUG ", log.LstdFlags),
			)
		}

		fn(cmd, args, c)
	}
}

func printError(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
