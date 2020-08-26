package reportsapi

import "github.com/lucassabreu/clockify-cli/http"

type FilterGroupEnum string

const ProjectFilterGroup = FilterGroupEnum("PROJECT")
const ClientFilterGroup = FilterGroupEnum("CLIENT")
const TaskFilterGroup = FilterGroupEnum("TASK")
const TagFilterGroup = FilterGroupEnum("TAG")
const DateFilterGroup = FilterGroupEnum("DATE")
const UserFilterGroup = FilterGroupEnum("USER")
const UserGroupFilterGroup = FilterGroupEnum("USER_GROUP")
const TimeEntryFilterGroup = FilterGroupEnum("TIMEENTRY")

type SummarySortColumnEnum string

const GroupSummarySortColumn = SummarySortColumnEnum("GROUP")
const DurationSummarySortColumn = SummarySortColumnEnum("DURATION")
const AmountSummarySortColumn = SummarySortColumnEnum("AMOUNT")

type SumarryFilter struct {
	Groups     []FilterGroupEnum
	SortColumn SummarySortColumnEnum
}

type SumarryRequest struct {
	DateRangeStart http.DateTime
	DateRangeEnd   http.DateTime

	SumarryFilter SumarryFilter
}
