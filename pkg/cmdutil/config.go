package cmdutil

import (
	"path"
	"strings"
	"time"

	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	CONF_WORKWEEK_DAYS       = "workweek-days"
	CONF_INTERACTIVE         = "interactive"
	CONF_ALLOW_NAME_FOR_ID   = "allow-name-for-id"
	CONF_USER_ID             = "user.id"
	CONF_WORKSPACE           = "workspace"
	CONF_TOKEN               = "token"
	CONF_ALLOW_INCOMPLETE    = "allow-incomplete"
	CONF_SHOW_TASKS          = "show-task"
	CONF_DESCR_AUTOCOMP      = "description-autocomplete"
	CONF_DESCR_AUTOCOMP_DAYS = "description-autocomplete-days"
	CONF_SHOW_TOTAL_DURATION = "show-total-duration"
	CONF_DEBUG               = "debug"
)

// Config manages configs and parameters used locally by the CLI
type Config interface {
	GetBool(string) bool
	SetBool(string, bool)

	GetInt(string) int
	SetInt(string, int)

	GetString(string) string
	SetString(string, string)

	GetStringSlice(string) []string
	SetStringSlice(string, []string)

	IsDebuging() bool
	IsAllowNameForID() bool
	IsInteractive() bool
	GetWorkWeekdays() []string

	Get(string) interface{}
	All() map[string]interface{}

	Save() error
}

type config struct{}

func (*config) GetBool(param string) bool {
	return viper.GetBool(param)
}

func (*config) SetBool(p string, b bool) {
	viper.Set(p, b)
}

func (*config) GetString(param string) string {
	return viper.GetString(param)
}

func (*config) SetString(p string, s string) {
	viper.Set(p, s)
}

func (*config) GetInt(param string) int {
	return viper.GetInt(param)
}

func (*config) SetInt(p string, i int) {
	viper.Set(p, i)
}

func (*config) GetStringSlice(param string) []string {
	return viper.GetStringSlice(param)
}

func (*config) SetStringSlice(p string, ss []string) {
	viper.Set(p, ss)
}

func (c *config) IsDebuging() bool {
	return c.GetBool(CONF_DEBUG)
}

func (c *config) GetWorkWeekdays() []string {
	return strhlp.Map(strings.ToLower, c.GetStringSlice(CONF_WORKWEEK_DAYS))
}

func (c *config) IsAllowNameForID() bool {
	return c.GetBool(CONF_ALLOW_NAME_FOR_ID)
}

func (c *config) IsInteractive() bool {
	return c.GetBool(CONF_INTERACTIVE)
}

func (*config) Get(p string) interface{} {
	return viper.Get(p)
}

func (*config) All() map[string]interface{} {
	return viper.AllSettings()
}

func (*config) Save() error {
	filename := viper.ConfigFileUsed()
	if filename == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		filename = path.Join(home, ".clockify-cli.yaml")
	}

	return viper.WriteConfigAs(filename)
}

func configFunc() func() (c Config) {
	return func() Config {
		return &config{}
	}
}

// GetWeekdays with their names
func GetWeekdays() []string {
	return []string{
		time.Sunday:    strings.ToLower(time.Sunday.String()),
		time.Monday:    strings.ToLower(time.Monday.String()),
		time.Tuesday:   strings.ToLower(time.Tuesday.String()),
		time.Wednesday: strings.ToLower(time.Wednesday.String()),
		time.Thursday:  strings.ToLower(time.Thursday.String()),
		time.Friday:    strings.ToLower(time.Friday.String()),
		time.Saturday:  strings.ToLower(time.Saturday.String()),
	}
}
