package reportsapi

import (
	"fmt"
	"time"

	"github.com/lucassabreu/clockify-cli/http"
)

// BaseUrl is the base url for all endpoints for "reporting" on Clockify
const BaseUrl = "https://reports.api.clockify.me/v1"

type ReportsClient struct {
	*http.Client
}

func NewReportsClient(apiKey string) (*ReportsClient, error) {
	c, err := http.NewHttpClient(BaseUrl, apiKey)
	if err != nil {
		return nil, err
	}

	return &ReportsClient{Client: c}, nil
}

func (c *ReportsClient) logf(format string, v ...interface{}) {
	if c.Logger == nil {
		return
	}
	c.Logger.Printf(format, v)
}

type TimeEntryParam struct {
	Workspace          string
	DateRangeStart     time.Time
	DateRangeEnd       time.Time
	Billable           *bool
	Description        string
	WithoutDescription bool
}

func (p TimeEntryParam) Fill(b BaseRequest) BaseRequest {
	b.DateRangeStart = http.NewDateTime(p.DateRangeStart)
	b.DateRangeEnd = http.NewDateTime(p.DateRangeEnd)
	b.AmountShown = AmountShownHideAmount
	return b
}

type entityStatus string

const (
	All      entityStatus = "ALL"
	Active   entityStatus = "ACTIVE"
	Inactive entityStatus = "INACTIVE"
	Archived entityStatus = "ARCHIVED"
	Done     entityStatus = "DONE"
)

type entityContains string

const (
	Contains       entityContains = "CONTAINS"
	DoesNotContain entityContains = "DOES_NOT_CONTAIN"
	ContainsOnly   entityContains = "CONTAINS_ONLY"
)

type EntityFilterParam struct {
	IDs      []string
	Status   entityStatus
	Contains entityContains
}

func (p EntityFilterParam) ToEntityFilter() EntityFilter {
	if p.IsEmpty() {
		return EntityFilter{}
	}

	e := EntityFilter{IDs: p.IDs}
	switch p.Status {
	case Active:
		e.Status = EntityFilterStatusActive
	case Inactive:
		e.Status = EntityFilterStatusInactive
	case Archived:
		e.Status = EntityFilterStatusArchived
	case Done:
		e.Status = EntityFilterStatusDone
	default:
		e.Status = EntityFilterStatusAll
	}

	return e
}

func (p EntityFilterParam) IsEmpty() bool {
	return len(p.IDs) == 0
}

type EntitiesParam struct {
	Users      EntityFilterParam
	UserGroups EntityFilterParam
	Clients    EntityFilterParam
	Projects   EntityFilterParam
	Tasks      EntityFilterParam
	Tags       EntityFilterParam
}

func (p EntitiesParam) Fill(b BaseRequest) BaseRequest {
	fn := func(r *EntityFilter, p EntityFilterParam) {
		if !p.IsEmpty() {
			*r = p.ToEntityFilter()
		}
	}

	fn(&b.Users, p.Users)
	fn(&b.UserGroups, p.UserGroups)
	fn(&b.Clients, p.Clients)
	fn(&b.Projects, p.Projects)
	fn(&b.Tasks, p.Tasks)
	fn(&b.Tags, p.Tags)

	return b
}

type sortOrder string

const (
	Ascending  sortOrder = "ASCENDING"
	Descending sortOrder = "DESCENDING"
)

func (s sortOrder) ToRequestSortOrder() SortOrder {
	switch s {
	case Descending:
		return SortOrderDescending
	case Ascending:
		return SortOrderAscending
	default:
		return SortOrderDefault
	}
}

type amountShown string

const (
	ShowAmount amountShown = "SHOW_AMOUNT"
	HideAmount amountShown = "HIDE_AMOUNT"
	Earned     amountShown = "EARNED"
	Cost       amountShown = "COST"
	Profit     amountShown = "PROFIT"
)

type SortAndAggregateParam struct {
	SortOrder  sortOrder
	AmoutShown amountShown
	Rouding    bool
}

type SummaryParam struct {
	TimeEntryParam
	EntitiesParam
}

func fill(filler ...interface{ Fill(BaseRequest) BaseRequest }) BaseRequest {
	b := BaseRequest{}
	for i := range filler {
		b = filler[i].Fill(b)
	}
	return b
}

func (c *ReportsClient) Summary(p SummaryParam) (SummaryReport, error) {
	var s SummaryReport

	r, err := c.NewRequest(
		"POST",
		fmt.Sprintf(
			"workspaces/%s/reports/summary",
			p.Workspace,
		),
		SummaryRequest{
			BaseRequest: fill(
				p.TimeEntryParam,
				p.EntitiesParam,
			),
			SummaryFilter: SummaryFilter{
				Groups: []FilterGroup{FilterGroupClient, FilterGroupProject},
			},
		},
	)
	if err != nil {
		return s, err
	}

	_, err = c.Do(r, &s)

	return s, err
}
