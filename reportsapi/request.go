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

type SummaryUsersContains string

const (
	SummaryUsersContainsDefault        SummaryUsersContains = ""
	SummaryUsersContainsContains       SummaryUsersContains = "CONTAINS"
	SummaryUsersContainsDoesNotContain SummaryUsersContains = "DOES_NOT_CONTAIN"
	SummaryUsersContainsContainsOnly   SummaryUsersContains = "CONTAINS_ONLY"
)

type SummaryUsersStatus string

const (
	SummaryUsersStatusAll      SummaryUsersStatus = "ALL"
	SummaryUsersStatusActive   SummaryUsersStatus = "ACTIVE"
	SummaryUsersStatusArchived SummaryUsersStatus = "ARCHIVED"
	SummaryUsersStatusInactive SummaryUsersStatus = "INACTIVE"
	SummaryUsersStatusDone     SummaryUsersStatus = "DONE"
)

type SummaryUsers struct {
	IDs      []string
	Contains SummaryUsersContains
	Status   SummaryUsersStatus
}

type SummaryRequest struct {
	DateRangeStart http.DateTime
	DateRangeEnd   http.DateTime
	SummaryFilter  SummaryFilter
	SortOrder      SortOrder
	ExportType     ExportType
	Rouding        bool
	AmountShown    AmountShown
	Users          SummaryUsers
}
