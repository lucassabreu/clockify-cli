package api

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// Client will help to access Clockify API
type Client struct {
	baseURL *url.URL
	http.Client
}

// baseURL is the Clockify API base URL
const baseURL = "https://api.clockify.me/api"

// ErrorMissingAPIKey returned if X-Api-Key is missing
var ErrorMissingAPIKey = errors.New("api Key must be informed")

// NewClient create a new Client
func NewClient(apiKey string) (*Client, error) {
	if len(apiKey) == 0 {
		return nil, ErrorMissingAPIKey
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		baseURL: u,
		Client: http.Client{
			Transport: transport{
				apiKey: apiKey,
				next:   http.DefaultTransport,
			},
		},
	}

	return c, nil
}

// WorkspacesFilter will be used to filter the workspaces
type WorkspacesFilter struct {
	Name string
}

// Workspaces list all the user's workspaces
func (c *Client) Workspaces(f WorkspacesFilter) ([]dto.Workspace, error) {
	var w []dto.Workspace

	r, err := c.NewRequest("GET", "workspaces/", nil)
	if err != nil {
		return w, err
	}

	_, err = c.Do(r, &w)

	if err != nil {
		return w, err
	}

	if f.Name == "" {
		return w, nil
	}

	ws := []dto.Workspace{}

	for _, i := range w {
		if strings.Contains(strings.ToLower(i.Name), strings.ToLower(f.Name)) {
			ws = append(ws, i)
		}
	}

	return ws, nil
}
