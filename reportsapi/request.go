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

type SummarySortColumn string

const (
	SummarySortColumnDefault  SummarySortColumn = ""
	SummarySortColumnGroup    SummarySortColumn = "GROUP"
	SummarySortColumnDuration SummarySortColumn = "DURATION"
	SummarySortColumnAmount   SummarySortColumn = "AMOUNT"
)

type SummaryFilter struct {
	Groups     []FilterGroup
	SortColumn SummarySortColumn
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

type SummaryRequest struct {
	DateRangeStart     http.DateTime
	DateRangeEnd       http.DateTime
	SummaryFilter      SummaryFilter
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
