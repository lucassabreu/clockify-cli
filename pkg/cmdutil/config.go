package cmdutil

import (
	"path"
	"strings"
	"sync"
	"time"

	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

const (
	CONF_WORKWEEK_DAYS                    = "workweek-days"
	CONF_INTERACTIVE                      = "interactive"
	CONF_ALLOW_NAME_FOR_ID                = "allow-name-for-id"
	CONF_SEARCH_PROJECTS_WITH_CLIENT_NAME = "search-project-with-client"
	CONF_USER_ID                          = "user.id"
	CONF_WORKSPACE                        = "workspace"
	CONF_TOKEN                            = "token"
	CONF_ALLOW_INCOMPLETE                 = "allow-incomplete"
	CONF_SHOW_TASKS                       = "show-task"
	CONF_SHOW_CLIENT                      = "show-client"
	CONF_DESCR_AUTOCOMP                   = "description-autocomplete"
	CONF_DESCR_AUTOCOMP_DAYS              = "description-autocomplete-days"
	CONF_SHOW_TOTAL_DURATION              = "show-total-duration"
	CONF_LOG_LEVEL                        = "log-level"
	CONF_ALLOW_ARCHIVED_TAGS              = "allow-archived-tags"
	CONF_INTERACTIVE_PAGE_SIZE            = "interactive-page-size"
	CONF_LANGUAGE                         = "lang"
	CONF_TIMEZONE                         = "time-zone"
)

const (
	LOG_LEVEL_NONE  = "none"
	LOG_LEVEL_DEBUG = "debug"
	LOG_LEVEL_INFO  = "info"
)

// Config manages configs and parameters used locally by the CLI
type Config interface {
	// GetBool retrieves a config by its name as a bool
	GetBool(string) bool
	// SetBool changes a bool config by its name
	SetBool(string, bool)

	// GetInt retrieves a config by its name as a int
	GetInt(string) int
	// SetInt changes a int config by its name
	SetInt(string, int)

	// GetString retrieves a config by its name as a string
	GetString(string) string
	// SetString changes a string config by its name
	SetString(string, string)

	// SetStringSlice retrieves a config by its name as a []string
	GetStringSlice(string) []string
	// SetStringSlice changes a []string config by its name
	SetStringSlice(string, []string)

	// IsDebuging configures CLI to log most of the data being used
	IsDebuging() bool
	// IsAllowNameForID configures the CLI to lookup entities ids by their name
	IsAllowNameForID() bool
	// IsInteractive configures the CLI to prompt the user interactively
	IsInteractive() bool
	// GetWorkWeekdays set which days of the week the user is expected to work
	GetWorkWeekdays() []string
	// InteractivePageSize sets how many items are shown when prompting
	// projects
	InteractivePageSize() int
	// IsSearchProjectWithClientsName defines if the project name for ID should
	// include the client's name
	IsSearchProjectWithClientsName() bool

	// Language what is the language to used when printing numbers
	Language() language.Tag

	// SetLanguage changes the language used for printing numbers
	SetLanguage(language.Tag)

	// TimeZone which time zone to use for showing date & time
	TimeZone() *time.Location

	// SetTimeZone changes the timezone used for dates
	SetTimeZone(*time.Location)

	// Get retrieves a config by its name
	Get(string) interface{}
	// All retrieves all the configurations of the CLI as a map
	All() map[string]interface{}

	// LogLevel sets how much should be logged during execution
	LogLevel() string

	// Save will persist the changes made to the configuration
	Save() error
}

type config struct {
	lngOnce sync.Once
	lang    language.Tag

	tzOnce   sync.Once
	timezone *time.Location
}

// IsSearchProjectWithClientsName defines if the project name for ID should
// include the client's name
func (c *config) IsSearchProjectWithClientsName() bool {
	return c.GetBool(CONF_SEARCH_PROJECTS_WITH_CLIENT_NAME)
}

func (c *config) InteractivePageSize() int {
	i := c.GetInt(CONF_INTERACTIVE_PAGE_SIZE)
	if i <= 0 {
		return 7
	}
	return i
}

func (c *config) LogLevel() string {
	l := c.GetString(CONF_LOG_LEVEL)
	switch l {
	case LOG_LEVEL_INFO, LOG_LEVEL_DEBUG:
		return l
	default:
		return LOG_LEVEL_NONE
	}
}

func (*config) GetBool(param string) bool {
	return viper.GetBool(param)
}

func (*config) SetBool(p string, b bool) {
	viper.Set(p, b)
}

func (*config) GetString(param string) string {
	return viper.GetString(param)
}

func (*config) SetString(p, s string) {
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
	return c.LogLevel() == LOG_LEVEL_DEBUG
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

func (c *config) SetLanguage(l language.Tag) {
	c.lngOnce = sync.Once{}
	c.SetString(CONF_LANGUAGE, l.String())
}

func (c *config) Language() language.Tag {
	c.lngOnce.Do(func() {
		c.lang = language.English

		lang := c.GetString(CONF_LANGUAGE)
		if lang == "" {
			return
		}

		if l, err := language.Parse(lang); err == nil {
			c.lang = l
			return
		}
	})

	return c.lang
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

// SetTimeZone changes the timezone used for dates
func (c *config) SetTimeZone(tz *time.Location) {
	c.tzOnce = sync.Once{}
	c.SetString(CONF_TIMEZONE, tz.String())
}

// TimeZone which time zone to use for showing date & time
func (c *config) TimeZone() *time.Location {
	c.tzOnce.Do(func() {
		c.timezone = time.Local

		tz := c.GetString(CONF_TIMEZONE)
		if tz == "" {
			return
		}

		if tz, err := time.LoadLocation(tz); err == nil {
			c.timezone = tz
			return
		}
	})

	return c.timezone
}
