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
	if len(p.IDs) == 0 {
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

type EntitiesParam struct {
	Users      EntityFilterParam
	UserGroups EntityFilterParam
	Clients    EntityFilterParam
	Projects   EntityFilterParam
	Tasks      EntityFilterParam
	Tags       EntityFilterParam
}

func (p EntitiesParam) Fill(b BaseRequest) BaseRequest {
	if !p.Users.IsEmpty() {

	}

	return b
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
