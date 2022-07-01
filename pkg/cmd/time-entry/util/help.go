package util

import "github.com/lucassabreu/clockify-cli/pkg/timeentryhlp"

const (
	HelpTimeEntryNowIfNotSet = "If no start time (`--when`) is set then the " +
		"current time will be used.\n"

	HelpInteractiveByDefault = "By default, the CLI will ask the " +
		"information interactively; use `--interactive=0` to disable it.\n" +
		"\n" +
		"If you prefer that it never don't do that by default, " +
		"run the bellow command, and use `--interactive` when you want " +
		"to be asked:\n" +
		"```\n" +
		"$ clockify-cli config set interactive false\n" +
		"```\n"

	HelpDateTimeFormats = "" +
		` - Full Date and Time:                "2016-02-01 15:04:05"` + "\n" +
		` - Date and Time (assumes 0 seconds): "2016-02-01 15:04"` + "\n" +
		` - Yesterday with Time:               "yesterday 15:04:05"` + "\n" +
		` - Yesterday with Time (0 seconds):   "yesterday 15:04"` + "\n" +
		` - Today at Time:                     "15:04:05"` + "\n" +
		` - Today at Time (assumes 0 seconds): "15:04"` + "\n" +
		` - 10mins in the future:              +10m` + "\n" +
		` - 1min and 30s ago:                  -90s` + "\n" +
		` - 1hour and 10min ago:               -1:10s` + "\n" +
		` - 1day, 10min and 30s ago:           -1d10m30s` + "\n"

	HelpTimeInputOnTimeEntry = "When setting a date/time input " +
		"(`--when` and `--when-to-close`) you can use any of the following " +
		"formats to set then:\n" +
		HelpDateTimeFormats

	HelpNamesForIds = "To be able to use names of resources instead of its " +
		"IDs you must enable the feature 'allow-name-for-id', to do that " +
		"run the command (the commands may take longer to look for the " +
		"resource id):\n" +
		"```\n" +
		"$ clockify-cli config set allow-name-for-id true\n" +
		"```\n\n"

	HelpValidateIncomplete = "By default, the CLI (and Clockify API) only " +
		"validates if the workspace and project rules are respected when a " +
		"time entry is stopped, if you prefer to validate when " +
		"starting/inserting it run the following command:\n" +
		"```\n" +
		"$ clockify-cli config set allow-incomplete false\n" +
		"```\n\n"

	HelpMoreInfoAboutStarting = "Use `clockify-cli in --help` for more " +
		"information about creating new time entries."

	HelpMoreInfoAboutPrinting = "Use `clockify-cli report --help` for more " +
		"information about printing time entries."

	HelpTimeEntriesAliasForEdit = "" +
		`If you want to edit the current (running) time entry you can ` +
		`use "` + timeentryhlp.AliasCurrent + `" instead of its ID.` + "\n" +
		`To edit the last ended time entry you can use "` +
		timeentryhlp.AliasLast + `" for it, for the one before that you ` +
		`can use "^2", for the previous "^3" and so on.` + "\n"
)
