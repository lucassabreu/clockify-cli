package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fullTimeFormat = "2006-01-02 15:04:05"
var simplerTimeFormat = "2006-01-02 15:04"
var onlyTimeFormat = "15:04:05"
var simplerOnlyTimeFormat = "15:04:05"

func withClockifyClient(fn func(cmd *cobra.Command, args []string, c *api.Client)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		c, err := api.NewClient(viper.GetString("token"))
		if err != nil {
			printError(err)
			return
		}

		if viper.GetBool("debug") {
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

func convertToTime(timeString string) (t time.Time, err error) {
	timeString = strings.TrimSpace(timeString)
	format := ""

	switch len(timeString) {
	case len(fullTimeFormat):
		format = fullTimeFormat
	case len(simplerTimeFormat):
		format = simplerTimeFormat
	case len(onlyTimeFormat):
		format = onlyTimeFormat
	case len(simplerOnlyTimeFormat):
		format = simplerOnlyTimeFormat
	}

	return time.ParseInLocation(format, timeString, time.Local)
}
