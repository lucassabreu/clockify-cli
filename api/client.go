package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// Client will help to access Clockify API
type Client struct {
	baseURL *url.URL
	http.Client
	debugLogger Logger
}

// baseURL is the Clockify API base URL
const baseURL = "https://api.clockify.me/api"

// ErrorMissingAPIKey returned if X-Api-Key is missing
var ErrorMissingAPIKey = errors.New("api Key must be informed")

// NewClient create a new Client, based on: https://clockify.github.io/clockify_api_docs/
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

// WorkspaceUsersParam params to query workspace users
type WorkspaceUsersParam struct {
	Workspace string
	Email     string
}

// WorkspaceUsers all users in a Workspace
func (c *Client) WorkspaceUsers(p WorkspaceUsersParam) ([]dto.User, error) {
	var users []dto.User

	r, err := c.NewRequest(
		"GET",
		fmt.Sprintf("workspaces/%s/users", p.Workspace),
		nil,
	)
	if err != nil {
		return users, err
	}

	_, err = c.Do(r, &users)

	if p.Email == "" {
		return users, nil
	}

	uCopy := []dto.User{}
	for _, i := range users {
		if strings.Contains(strings.ToLower(i.Email), strings.ToLower(p.Email)) {
			uCopy = append(uCopy, i)
		}
	}

	return uCopy, nil
}

// LogParam params to query entries
type LogParam struct {
	Workspace string
	UserID    string
	Date      time.Time
	AllPages  bool
}

// Log list time entries
func (c *Client) Log(p LogParam) ([]dto.TimeEntry, error) {
	c.debugf("Log - Date Param: %s", p.Date)

	var timeEntries []dto.TimeEntry

	d := p.Date.Round(time.Hour)
	d = d.Add(time.Hour * time.Duration(d.Hour()) * -1)

	filter := dto.TimeEntryStartEndRequest{
		Start: dto.DateTime{Time: d},
		End:   dto.DateTime{Time: d.Add(time.Hour * 24)},
	}

	c.debugf("Log Filter Params: Start: %s, End: %s", filter.Start, filter.End)

	r, err := c.NewRequest(
		"POST",
		fmt.Sprintf(
			"workspaces/%s/timeEntries/user/%s/entriesInRange",
			p.Workspace,
			p.UserID,
		),
		filter,
	)
	if err != nil {
		return timeEntries, err
	}

	_, err = c.Do(r, &timeEntries)

	return timeEntries, nil
}

// LogInProgressParam params to query entries
type LogInProgressParam struct {
	Workspace string
}

// LogInProgress show time entry in progress (if any)
func (c *Client) LogInProgress(p LogInProgressParam) (*dto.TimeEntryImpl, error) {
	var timeEntryImpl *dto.TimeEntryImpl

	r, err := c.NewRequest(
		"GET",
		fmt.Sprintf(
			"workspaces/%s/timeEntries/inProgress",
			p.Workspace,
		),
		nil,
	)

	if err != nil {
		return timeEntryImpl, err
	}

	_, err = c.Do(r, &timeEntryImpl)
	return timeEntryImpl, err
}

// GetTagParam params to find a tag
type GetTagParam struct {
	Workspace string
	TagID     string
}

// GetTag get a single tag, if it exists
func (c *Client) GetTag(p GetTagParam) (*dto.Tag, error) {
	var tag *dto.Tag

	r, err := c.NewRequest(
		"GET",
		fmt.Sprintf(
			"workspaces/%s/tags/%s",
			p.Workspace,
			p.TagID,
		),
		nil,
	)

	if err != nil {
		return tag, err
	}

	_, err = c.Do(r, &tag)
	return tag, err
}

// GetProjectParam params to get a Project
type GetProjectParam struct {
	Workspace string
	ProjectID string
}

// GetProject get a single Project, if exists
func (c *Client) GetProject(p GetProjectParam) (*dto.Project, error) {
	var project *dto.Project

	r, err := c.NewRequest(
		"GET",
		fmt.Sprintf(
			"workspaces/%s/projects/%s",
			p.Workspace,
			p.ProjectID,
		),
		nil,
	)

	if err != nil {
		return project, err
	}

	_, err = c.Do(r, &project)
	return project, err
}

// GetUser get a specific user by its id
func (c *Client) GetUser(id string) (*dto.User, error) {
	var user *dto.User

	r, err := c.NewRequest(
		"GET",
		fmt.Sprintf("users/%s", id),
		nil,
	)

	if err != nil {
		return user, err
	}

	_, err = c.Do(r, &user)
	return user, err
}

// GetTaskParam params to get a Task
type GetTaskParam struct {
	Workspace string
	TaskID    string
}

// GetTask get a single Task, if exists
func (c *Client) GetTask(p GetTaskParam) (*dto.Task, error) {
	var task *dto.Task

	r, err := c.NewRequest(
		"GET",
		fmt.Sprintf(
			"workspaces/%s/tasks/%s",
			p.Workspace,
			p.TaskID,
		),
		nil,
	)

	if err != nil {
		return task, err
	}

	_, err = c.Do(r, &task)
	return task, err
}
