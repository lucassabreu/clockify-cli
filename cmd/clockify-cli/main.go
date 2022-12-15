package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/lucassabreu/clockify-cli/pkg/cmd"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const (
	exitOK     = 0
	exitError  = 1
	exitCancel = 2
)

func main() {
	exitCode := execute()
	os.Exit(exitCode)
}

func execute() int {
	f := cmdutil.NewFactory(cmdutil.Version{
		Tag:    version,
		Commit: commit,
		Date:   date,
	})

	rootCmd := cmd.NewCmdRoot(f)
	rootCmd.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return cmdutil.FlagErrorWrap(err)
	})

	cmd := rootCmd
	err := bindViper(rootCmd)

	if err == nil {
		cmd, err = rootCmd.ExecuteC()
	}

	if err == nil {
		return exitOK
	}

	stderr := cmd.ErrOrStderr()
	if errors.Is(err, terminal.InterruptErr) {
		fmt.Fprintln(stderr)
		return exitCancel
	}

	var flagError *cmdutil.FlagError
	if errors.As(err, &flagError) {
		fmt.Fprintln(stderr, flagError.Error())
		fmt.Fprintln(stderr, cmd.UsageString())
		return exitError
	}

	if f.Config().IsDebuging() {
		fmt.Fprintf(stderr, "%+v\n", err)
	} else {
		fmt.Fprintln(stderr, err.Error())
	}

	return exitError
}

func bindViper(rootCmd *cobra.Command) error {
	envPrefix := "CLOCKIFY"
	bind := func(flag *pflag.Flag, conf, sufix string) error {
		if flag == nil {
			return nil
		}

		flag.Usage = flag.Usage +
			" (defaults to env $" + envPrefix + "_" + sufix + ")"

		return viper.BindPFlag(conf, flag)
	}

	var err error
	l := rootCmd.PersistentFlags().Lookup
	if err = bind(l("token"), cmdutil.CONF_TOKEN, "TOKEN"); err != nil {
		return err
	}

	err = bind(l("workspace"), cmdutil.CONF_WORKSPACE, "WORKSPACE")
	if err != nil {
		return err
	}

	if err = bind(l("user-id"), cmdutil.CONF_USER_ID, "USER_ID"); err != nil {
		return err
	}

	err = bind(l("log-level"), cmdutil.CONF_LOG_LEVEL, "LOG_LEVEL")
	if err != nil {
		return err
	}

	viper.RegisterAlias(cmdutil.CONF_ALLOW_NAME_FOR_ID, "allow-project-name")
	if err = bind(l("allow-name-for-id"), cmdutil.CONF_ALLOW_NAME_FOR_ID,
		"ALLOW_NAME_FOR_ID"); err != nil {
		return err
	}

	if err = bind(l("interactive"), cmdutil.CONF_INTERACTIVE,
		"INTERACTIVE"); err != nil {
		return err
	}

	if err = bind(l("interactive-page-size"),
		cmdutil.CONF_INTERACTIVE_PAGE_SIZE,
		"INTERACTIVE_PAGE_SIZE"); err != nil {
		return err
	}

	f := l("interactive")
	f.Usage = f.Usage + "\n" +
		"You can be disable it temporally by setting it to 0 " +
		"(-i=0 or " + envPrefix + "_INTERACTIVE=0)"

	var cfgFile = ""
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.clockify-cli.yaml)")

	var viperErr error
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if viperErr != nil {
			return viperErr
		}

		if withTotals := cmd.Flags().Lookup("with-totals"); withTotals != nil {
			viper.SetDefault(cmdutil.CONF_SHOW_TOTAL_DURATION, true)
			if err := viper.BindPFlag(
				cmdutil.CONF_SHOW_TOTAL_DURATION, withTotals); err != nil {
				return err
			}
		}

		if flag := cmd.Flags().Lookup("allow-incomplete"); flag != nil {
			if err := bind(flag, cmdutil.CONF_ALLOW_INCOMPLETE,
				"ALLOW_INCOMPLETE"); err != nil {
				return err
			}
		}

		return nil
	}

	cobra.OnInitialize(func() {
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			home, err := homedir.Dir()
			if err != nil {
				viperErr = err
				return
			}

			viper.AddConfigPath(home)
			viper.SetConfigName(".clockify-cli")
		}

		viper.SetEnvPrefix(envPrefix)
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()

		err := viper.ReadInConfig()
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return
		}

		viperErr = err
	})

	return nil
}
