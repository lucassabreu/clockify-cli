package reportsapi

import "github.com/lucassabreu/clockify-cli/http"

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
	Groups     []FilterGroup
	SortColumn SortColumn
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

type SummaryEntityFilterContains string

const (
	SummaryEntityFilterContainsDefault        SummaryEntityFilterContains = ""
	SummaryEntityFilterContainsContains       SummaryEntityFilterContains = "CONTAINS"
	SummaryEntityFilterContainsDoesNotContain SummaryEntityFilterContains = "DOES_NOT_CONTAIN"
	SummaryEntityFilterContainsContainsOnly   SummaryEntityFilterContains = "CONTAINS_ONLY"
)

type SummaryEntityFilterStatus string

const (
	SummaryEntityFilterStatusAll      SummaryEntityFilterStatus = "ALL"
	SummaryEntityFilterStatusActive   SummaryEntityFilterStatus = "ACTIVE"
	SummaryEntityFilterStatusArchived SummaryEntityFilterStatus = "ARCHIVED"
	SummaryEntityFilterStatusInactive SummaryEntityFilterStatus = "INACTIVE"
	SummaryEntityFilterStatusDone     SummaryEntityFilterStatus = "DONE"
)

type SummaryEntityFilter struct {
	IDs      []string
	Contains SummaryEntityFilterContains
	Status   SummaryEntityFilterStatus
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

type SummaryCustomField struct {
	ID              string
	Value           string
	Type            CustomFieldType
	NumberCondition CustomTypeNumberCondition
	Empty           bool
}

type BaseRequest struct {
	DateRangeStart     http.DateTime
	DateRangeEnd       http.DateTime
	SortOrder          SortOrder
	ExportType         ExportType
	Rouding            bool
	AmountShown        AmountShown
	Users              SummaryEntityFilter
	UserGroups         SummaryEntityFilter
	Clients            SummaryEntityFilter
	Projects           SummaryEntityFilter
	Tasks              SummaryEntityFilter
	Tags               SummaryEntityFilter
	Billable           *bool
	Description        string
	WithoutDescription bool
	CustomFields       []SummaryCustomField
}

type SummaryRequest struct {
	SummaryFilter SummaryFilter
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

type SummaryResponse struct {
	Totals   []Total
	GroupOne []SummaryGroup
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

type DetailedResponse struct {
	Totals      []Total
	TimeEntries []DetailedTimeEntry
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

type DayTotal struct {
	Date     http.DateTime
	Amount   float32
	Duration int
}

type WeeklyResponse struct {
	Totals      []Total
	TotalsByDay []DayTotal
}
