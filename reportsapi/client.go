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

type SummaryParam struct {
	TimeEntryParam
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
			BaseRequest: BaseRequest{
				DateRangeStart: http.NewDateTime(p.DateRangeStart),
				DateRangeEnd:   http.NewDateTime(p.DateRangeEnd),
				AmountShown:    AmountShownHideAmount,
			},
			SummaryFilter: SummaryFilter{
				Groups: []FilterGroup{FilterGroupProject},
			},
		},
	)
	if err != nil {
		return s, err
	}

	_, err = c.Do(r, &s)

	return s, err
}
