package reportsapi

import (
	"strconv"
	"strings"
	"time"

	"github.com/lucassabreu/clockify-cli/http"
)

type FilterGroup string

const (
	FilterGroupProject   FilterGroup = "PROJECT"
	FilterGroupClient    FilterGroup = "CLIENT"
	FilterGroupTask      FilterGroup = "TASK"
	FilterGroupTag       FilterGroup = "TAG"
	FilterGroupDate      FilterGroup = "DATE"
	FilterGroupUser      FilterGroup = "USER"
	FilterGroupUserGroup FilterGroup = "USER_GROUP"
	FilterGroupTimeEntry FilterGroup = "TIMEENTRY"
)

type SortColumn string

const (
	SortColumnDefault  SortColumn = ""
	SortColumnGroup    SortColumn = "GROUP"
	SortColumnDuration SortColumn = "DURATION"
	SortColumnAmount   SortColumn = "AMOUNT"
)

type SummaryFilter struct {
	Groups     []FilterGroup `json:"groups"`
	SortColumn SortColumn    `json:"sortColumn,omitempty"`
}

type SortOrder string

const (
	SortOrderDefault    SortOrder = ""
	SortOrderAscending  SortOrder = "ASCENDING"
	SortOrderDescending SortOrder = "DESCENDING"
)

type ExportType string

const (
	ExportTypeDefault ExportType = ""
	ExportTypeJSON    ExportType = "JSON"
	ExportTypeCSV     ExportType = "CSV"
	ExportTypeXLSX    ExportType = "XLSX"
	ExportTypePDF     ExportType = "PDF"
)

type AmountShown string

const (
	AmountShownDefault    AmountShown = ""
	AmountShownShowAmount AmountShown = "SHOW_AMOUNT"
	AmountShownHideAmount AmountShown = "HIDE_AMOUNT"
	AmountShownEarned     AmountShown = "EARNED"
	AmountShownCost       AmountShown = "COST"
	AmountShownProfit     AmountShown = "PROFIT"
)

type EntityFilterContains string

const (
	EntityFilterContainsDefault        EntityFilterContains = ""
	EntityFilterContainsContains       EntityFilterContains = "CONTAINS"
	EntityFilterContainsDoesNotContain EntityFilterContains = "DOES_NOT_CONTAIN"
	EntityFilterContainsContainsOnly   EntityFilterContains = "CONTAINS_ONLY"
)

type EntityFilterStatus string

const (
	EntityFilterStatusAll      EntityFilterStatus = "ALL"
	EntityFilterStatusActive   EntityFilterStatus = "ACTIVE"
	EntityFilterStatusArchived EntityFilterStatus = "ARCHIVED"
	EntityFilterStatusInactive EntityFilterStatus = "INACTIVE"
	EntityFilterStatusDone     EntityFilterStatus = "DONE"
)

type EntityFilter struct {
	IDs      []string
	Contains EntityFilterContains
	Status   EntityFilterStatus
}

func mapString(s []string, f func(string) string) []string {
	ns := make([]string, len(s))
	for i := range s {
		ns[i] = f(s[i])
	}
	return ns
}

func quoteJoin(s []string) string {
	return strings.Join(
		mapString(s, strconv.Quote),
		",",
	)
}

// MarshalJSON converts DateTime correctly
func (e EntityFilter) MarshalJSON() ([]byte, error) {
	if len(e.IDs) == 0 {
		return []byte("null"), nil
	}

	if e.Contains == EntityFilterContainsDefault {
		e.Contains = EntityFilterContainsContains
	}

	if e.Status == EntityFilterStatus("") {
		e.Status = EntityFilterStatusAll
	}

	b := []byte(
		"{" +
			`"ids":[` + quoteJoin(e.IDs) + "]," +
			`"contains":"` + string(e.Contains) + `",` +
			`"status":"` + string(e.Status) + `"` +
			"}",
	)

	return b, nil
}

type CustomFieldType string

const (
	CustomFieldTypeTXT              CustomFieldType = "TXT"
	CustomFieldTypeNumber           CustomFieldType = "NUMBER"
	CustomFieldTypeDropdownSingle   CustomFieldType = "DROPDOWN_SINGLE"
	CustomFieldTypeDropdownMultiple CustomFieldType = "DROPDOWN_MULTIPLE"
	CustomFieldTypeCheckbox         CustomFieldType = "CHECKBOX"
	CustomFieldTypeLink             CustomFieldType = "LINK"
)

type CustomTypeNumberCondition string

type CustomFieldFilter struct {
	ID              string
	Value           string
	Type            CustomFieldType
	NumberCondition CustomTypeNumberCondition
	Empty           bool
}

type BaseRequest struct {
	DateRangeStart     http.DateTime       `json:"dateRangeStart"`
	DateRangeEnd       http.DateTime       `json:"dateRangeEnd"`
	SortOrder          SortOrder           `json:"sortOrder,omitempty"`
	ExportType         ExportType          `json:"exportType,omitempty"`
	Rouding            bool                `json:"rounding,omitempty"`
	AmountShown        AmountShown         `json:"amountShown,omitempty"`
	Users              EntityFilter        `json:"users,omitempty"`
	UserGroups         EntityFilter        `json:"userGroups,omitempty"`
	Clients            EntityFilter        `json:"clients,omitempty"`
	Projects           EntityFilter        `json:"projects,omitempty"`
	Tasks              EntityFilter        `json:"tasks,omitempty"`
	Tags               EntityFilter        `json:"tags,omitempty"`
	Billable           *bool               `json:"billable,omitempty"`
	Description        string              `json:"description,omitempty"`
	WithoutDescription bool                `json:"withoutDescription,omitempty"`
	CustomFields       []CustomFieldFilter `json:"customFields,omitempty"`
}

type SummaryRequest struct {
	SummaryFilter SummaryFilter `json:"summaryFilter"`
	BaseRequest
}

type AuditFilter struct {
	WithoutProject  bool
	WithoutTask     bool
	Duration        int
	DurationShorter bool
}

type DetailedFilter struct {
	Page        int
	PageSize    int
	SortColumn  SortColumn
	AuditFilter AuditFilter
}

type DetailedRequest struct {
	DetailedFilter DetailedFilter
	BaseRequest
}

type TimeInterval struct {
	Start    http.DateTime
	End      http.DateTime
	Duration int
}

type CustomField struct {
	CustomFieldID string
	TimeEntryID   string
	Value         string
	Name          string
}

type WeeklyFilterGroup string

const (
	WeeklyFilterGroupProject WeeklyFilterGroup = "PROJECT"
	WeeklyFilterGroupUser    WeeklyFilterGroup = "USER"
)

type WeeklyFilterSubGroup string

const (
	WeeklyFilterSubGroupTIME   WeeklyFilterSubGroup = "TIME"
	WeeklyFilterSubGroupEARNED WeeklyFilterSubGroup = "EARNED"
)

type WeeklyFilter struct {
	Group    WeeklyFilterGroup
	SubGroup string
}

type WeeklyRequest struct {
	WeeklyFilter
	BaseRequest
}

type Total struct {
	TotalTime         int
	TotalBillableTime int
	EntriesCount      int
	TotalAmount       float32
}

type SummaryGroup struct {
	ID       string
	Duration int
	Amount   float32
	Name     string
	Children []SummaryGroup
}

type SummaryReport struct {
	Totals   []Total
	GroupOne []SummaryGroup
}

type DetailedTimeEntry struct {
	ID           string
	Description  string
	UserID       string
	Billable     bool
	TaskID       *string
	ProjectID    *string
	TimeInterval TimeInterval
	Tags         []string
	IsLocked     bool
	CustomFields []CustomField
	Amount       float32
	Rate         float32
	UserName     string
	UserEmail    string
}

type DetailedReport struct {
	Totals      []Total
	TimeEntries []DetailedTimeEntry
}

type DayTotal struct {
	Date     time.Time
	Amount   float32
	Duration int
}

type WeeklyReport struct {
	Totals      []Total
	TotalsByDay []DayTotal
}
