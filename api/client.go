package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

func (c *Client) CurrentUser() (dto.User, error) {
	var u dto.User
	r, err := c.NewRequest("GET", "v1/user", nil)
	if err != nil {
		return u, err
	}

	_, err = c.Do(r, &u)
	return u, err
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
	if err != nil {
		return users, err
	}

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

// Log list time entries from a date
func (c *Client) Log(p LogParam) ([]dto.TimeEntry, error) {
	c.debugf("Log - Date Param: %s", p.Date)

	d := p.Date.Round(time.Hour)
	d = d.Add(time.Hour * time.Duration(d.Hour()) * -1)

	return c.LogRange(LogRangeParam{
		Workspace: p.Workspace,
		UserID:    p.UserID,
		AllPages:  p.AllPages,
		FirstDate: d,
		LastDate:  d.Add(time.Hour * 24),
	})
}

// LogRangeParam params to query entries
type LogRangeParam struct {
	Workspace string
	UserID    string
	FirstDate time.Time
	LastDate  time.Time
	AllPages  bool
}

// LogRange list time entries by date range
func (c *Client) LogRange(p LogRangeParam) ([]dto.TimeEntry, error) {
	c.debugf("LogRange - First Date Param: %s | Last Date Param: %s", p.FirstDate, p.LastDate)

	var timeEntries []dto.TimeEntry

	filter := dto.TimeEntryStartEndRequest{
		Start: dto.DateTime{Time: p.FirstDate},
		End:   dto.DateTime{Time: p.LastDate},
	}

	c.debugf("Log Filter Params: Start: %s, End: %s", filter.Start, filter.End)

	r, err := c.NewRequest(
		"GET",
		fmt.Sprintf(
			"v1/workspaces/%s/user/%s/time-entries",
			p.Workspace,
			p.UserID,
		),
		filter,
	)
	if err != nil {
		return timeEntries, err
	}

	_, err = c.Do(r, &timeEntries)

	return timeEntries, err
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

// CreateTimeEntryParam params to create a new time entry
type CreateTimeEntryParam struct {
	Workspace   string
	Start       time.Time
	End         *time.Time
	Billable    bool
	Description string
	ProjectID   string
	TaskID      string
	TagIDs      []string
}

// CreateTimeEntry create a new time entry
func (c *Client) CreateTimeEntry(p CreateTimeEntryParam) (dto.TimeEntryImpl, error) {
	var t dto.TimeEntryImpl

	var end *dto.DateTime
	if p.End != nil {
		end = &dto.DateTime{Time: *p.End}
	}

	r, err := c.NewRequest(
		"POST",
		fmt.Sprintf(
			"workspaces/%s/timeEntries/",
			p.Workspace,
		),
		dto.CreateTimeEntryRequest{
			Start:       dto.DateTime{Time: p.Start},
			End:         end,
			Billable:    p.Billable,
			Description: p.Description,
			ProjectID:   p.ProjectID,
			TaskID:      p.TaskID,
			TagIDs:      p.TagIDs,
		},
	)

	if err != nil {
		return t, err
	}

	_, err = c.Do(r, &t)

	return t, err
}

// GetTagsParam params to get all tags of a workspace
type GetTagsParam struct {
	Workspace string
}

// GetTags get all tags of a workspace
func (c *Client) GetTags(p GetTagsParam) ([]dto.Tag, error) {
	var ps []dto.Tag

	r, err := c.NewRequest(
		"GET",
		fmt.Sprintf(
			"workspaces/%s/tags",
			p.Workspace,
		),
		nil,
	)

	if err != nil {
		return ps, err
	}

	_, err = c.Do(r, &ps)
	return ps, err
}

// GetProjectsParam params to get all project of a workspace
type GetProjectsParam struct {
	Workspace string
}

// GetProjects get all project of a workspace
func (c *Client) GetProjects(p GetProjectsParam) ([]dto.Project, error) {
	var ps []dto.Project

	r, err := c.NewRequest(
		"GET",
		fmt.Sprintf(
			"workspaces/%s/projects/",
			p.Workspace,
		),
		nil,
	)

	if err != nil {
		return ps, err
	}

	_, err = c.Do(r, &ps)
	return ps, err
}

// OutParam params to end the current time entry
type OutParam struct {
	Workspace string
	End       time.Time
}

// Out create a new time entry
func (c *Client) Out(p OutParam) error {
	r, err := c.NewRequest(
		"PUT",
		fmt.Sprintf(
			"workspaces/%s/timeEntries/endStarted",
			p.Workspace,
		),
		dto.OutTimeEntryRequest{
			End: dto.DateTime{Time: p.End},
		},
	)

	if err != nil {
		return err
	}

	_, err = c.Do(r, nil)
	return err
}

// UpdateTimeEntryParam params to update a new time entry
type UpdateTimeEntryParam struct {
	Workspace   string
	TimeEntryID string
	Start       time.Time
	End         *time.Time
	Billable    bool
	Description string
	ProjectID   string
	TaskID      string
	TagIDs      []string
}

// UpdateTimeEntry update a time entry
func (c *Client) UpdateTimeEntry(p UpdateTimeEntryParam) (dto.TimeEntryImpl, error) {
	var t dto.TimeEntryImpl

	var end *dto.DateTime
	if p.End != nil {
		end = &dto.DateTime{Time: *p.End}
	}

	r, err := c.NewRequest(
		"PUT",
		fmt.Sprintf(
			"workspaces/%s/timeEntries/%s",
			p.Workspace,
			p.TimeEntryID,
		),
		dto.UpdateTimeEntryRequest{
			Start:       dto.DateTime{Time: p.Start},
			End:         end,
			Billable:    p.Billable,
			Description: p.Description,
			ProjectID:   p.ProjectID,
			TaskID:      p.TaskID,
			TagIDs:      p.TagIDs,
		},
	)

	if err != nil {
		return t, err
	}

	_, err = c.Do(r, &t)

	return t, err
}

// GetRecentTimeEntries params to get recent time entries
type GetRecentTimeEntries struct {
	Workspace    string
	UserID       string
	Page         int
	ItemsPerPage int
}

// GetRecentTimeEntries will return the recent time entries of the user, paginated
func (c *Client) GetRecentTimeEntries(p GetRecentTimeEntries) (dto.TimeEntriesList, error) {
	var resp dto.TimeEntriesList

	r, err := c.NewRequest(
		"GET",
		fmt.Sprintf(
			"workspaces/%s/timeEntries/user/%s",
			p.Workspace,
			p.UserID,
		),
		nil,
	)

	if p.Page != 0 {
		r.URL.Query().Add("page", strconv.Itoa(p.Page))
	}

	if p.ItemsPerPage != 0 {
		r.URL.Query().Add("limit", strconv.Itoa(p.ItemsPerPage))
	}

	if err != nil {
		return resp, err
	}

	_, err = c.Do(r, &resp)
	return resp, err
}
