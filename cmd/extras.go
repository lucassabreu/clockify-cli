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
var simplerOnlyTimeFormat = "15:04"

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
	if len(fullTimeFormat) != len(timeString) && len(simplerTimeFormat) != len(timeString) && len(onlyTimeFormat) != len(timeString) && len(simplerOnlyTimeFormat) != len(timeString) {
		return t, fmt.Errorf(
			"supported formats are: %s",
			strings.Join(
				[]string{fullTimeFormat, simplerTimeFormat, onlyTimeFormat, simplerOnlyTimeFormat},
				", ",
			),
		)
	}

	if len(simplerOnlyTimeFormat) == len(timeString) || len(simplerTimeFormat) == len(timeString) {
		timeString = timeString + ":00"
	}

	if len(onlyTimeFormat) == len(timeString) {
		timeString = time.Now().Format("2006-01-02") + " " + timeString
	}

	return time.ParseInLocation(fullTimeFormat, timeString, time.Local)
}
