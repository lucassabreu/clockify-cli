package mocks

import (
	"time"

	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"golang.org/x/text/language"
)

// SimpleConfig is used to set configs for tests were changing the config or
// accessing them with Get and All is not important
type SimpleConfig struct {
	WorkweekDays                 []string
	Interactive                  bool
	InteractivePageSizeNumber    int
	AllowNameForID               bool
	UserID                       string
	Workspace                    string
	Token                        string
	AllowIncomplete              bool
	ShowTask                     bool
	DescriptionAutocomplete      bool
	DescriptionAutocompleteDays  int
	ShowTotalDuration            bool
	LogLevelValue                string
	AllowArchivedTags            bool
	SearchProjectWithClientsName bool
	LanguageTag                  language.Tag
	TimeZoneLoc                  *time.Location
}

func (d *SimpleConfig) GetBool(n string) bool {
	switch n {
	case cmdutil.CONF_INTERACTIVE:
		return d.Interactive
	case cmdutil.CONF_ALLOW_NAME_FOR_ID:
		return d.AllowNameForID
	case cmdutil.CONF_ALLOW_INCOMPLETE:
		return d.AllowIncomplete
	case cmdutil.CONF_SHOW_TASKS:
		return d.ShowTask
	case cmdutil.CONF_DESCR_AUTOCOMP:
		return d.DescriptionAutocomplete
	case cmdutil.CONF_SHOW_TOTAL_DURATION:
		return d.ShowTotalDuration
	case cmdutil.CONF_ALLOW_ARCHIVED_TAGS:
		return d.AllowArchivedTags
	default:
		return false
	}
}

func (*SimpleConfig) SetBool(_ string, _ bool) {
	panic("should not call")
}

func (d *SimpleConfig) GetInt(n string) int {
	switch n {
	case cmdutil.CONF_DESCR_AUTOCOMP_DAYS:
		return d.DescriptionAutocompleteDays
	case cmdutil.CONF_INTERACTIVE_PAGE_SIZE:
		return d.InteractivePageSize()
	default:
		return 0
	}
}

func (*SimpleConfig) SetInt(_ string, _ int) {
	panic("should not call")
}

func (d *SimpleConfig) GetString(n string) string {
	switch n {
	case cmdutil.CONF_USER_ID:
		return d.UserID
	case cmdutil.CONF_WORKSPACE:
		return d.Workspace
	case cmdutil.CONF_TOKEN:
		return d.Token
	case cmdutil.CONF_LOG_LEVEL:
		return d.LogLevelValue
	default:
		return ""

	}
}

func (*SimpleConfig) SetString(_, _ string) {
	panic("should not call")
}

func (d *SimpleConfig) GetStringSlice(n string) []string {
	switch n {
	case cmdutil.CONF_WORKWEEK_DAYS:
		return d.WorkweekDays
	default:
		return []string{}
	}
}

func (*SimpleConfig) SetStringSlice(_ string, _ []string) {
	panic("should not call")
}

func (d *SimpleConfig) IsDebuging() bool {
	return d.LogLevel() == cmdutil.LOG_LEVEL_DEBUG
}

func (d *SimpleConfig) IsAllowNameForID() bool {
	return d.AllowNameForID
}

func (d *SimpleConfig) IsInteractive() bool {
	return d.Interactive
}

func (d *SimpleConfig) GetWorkWeekdays() []string {
	return d.WorkweekDays
}

func (d *SimpleConfig) SetLanguage(l language.Tag) {
	d.LanguageTag = l
}

func (d *SimpleConfig) Language() language.Tag {
	return d.LanguageTag
}

func (*SimpleConfig) Get(_ string) interface{} {
	panic("should not call")
}

func (*SimpleConfig) All() map[string]interface{} {
	panic("should not call")
}

func (d *SimpleConfig) LogLevel() string {
	return d.LogLevelValue
}

// TimeZone which time zone to use for showing date & time
func (s *SimpleConfig) TimeZone() *time.Location {
	if s.TimeZoneLoc == nil {
		s.TimeZoneLoc = time.UTC
	}

	return s.TimeZoneLoc
}

// SetTimeZone changes the timezone used for dates
func (s *SimpleConfig) SetTimeZone(loc *time.Location) {
	s.TimeZoneLoc = loc
}

// IsSearchProjectWithClientsName defines if the project name for ID should
// include the client's name
func (s *SimpleConfig) IsSearchProjectWithClientsName() bool {
	return s.SearchProjectWithClientsName
}

// InteractivePageSize sets how many items are shown when prompting
// projects
func (s *SimpleConfig) InteractivePageSize() int {
	return s.InteractivePageSizeNumber
}

func (*SimpleConfig) Save() error {
	panic("should not call")
}
