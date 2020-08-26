package reportsapi

import "github.com/lucassabreu/clockify-cli/http"

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
