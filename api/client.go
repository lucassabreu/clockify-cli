package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
	stackedErrors "github.com/pkg/errors"
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
		return nil, stackedErrors.WithStack(ErrorMissingAPIKey)
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, stackedErrors.WithStack(err)
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

// GetWorkspaces will be used to filter the workspaces
type GetWorkspaces struct {
	Name string
}

// Workspaces list all the user's workspaces
func (c *Client) GetWorkspaces(f GetWorkspaces) ([]dto.Workspace, error) {
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

// PaginationParam parameters about pagination
type PaginationParam struct {
	AllPages bool
	Page     int
	PageSize int
}

// LogParam params to query entries
type LogParam struct {
	Workspace string
	UserID    string
	Date      time.Time
	PaginationParam
}

// Log list time entries from a date
func (c *Client) Log(p LogParam) ([]dto.TimeEntry, error) {
	c.debugf("Log - Date Param: %s", p.Date)

	d := p.Date.Round(time.Hour)
	d = d.Add(time.Hour * time.Duration(d.Hour()) * -1)

	return c.LogRange(LogRangeParam{
		Workspace:       p.Workspace,
		UserID:          p.UserID,
		FirstDate:       d,
		LastDate:        d.Add(time.Hour * 24),
		PaginationParam: p.PaginationParam,
	})
}

// LogRangeParam params to query entries
type LogRangeParam struct {
	Workspace string
	UserID    string
	FirstDate time.Time
	LastDate  time.Time
	PaginationParam
}

// LogRange list time entries by date range
func (c *Client) LogRange(p LogRangeParam) ([]dto.TimeEntry, error) {
	c.debugf("LogRange - First Date Param: %s | Last Date Param: %s", p.FirstDate, p.LastDate)

	var timeEntries []dto.TimeEntry

	b := true
	filter := dto.TimeEntryStartEndRequest{
		Start:    dto.DateTime{Time: p.FirstDate},
		End:      dto.DateTime{Time: p.LastDate},
		Hydrated: &b,
	}

	c.debugf("Log Filter Params: Start: %s, End: %s", filter.Start, filter.End)

	var tes []dto.TimeEntry
	err := c.paginate(
		"GET",
		fmt.Sprintf(
			"v1/workspaces/%s/user/%s/time-entries",
			p.Workspace,
			p.UserID,
		),
		p.PaginationParam,
		filter,
		&tes,
		func(res interface{}) (int, error) {
			if res == nil {
				return 0, nil
			}

			tes := res.(*[]dto.TimeEntry)
			timeEntries = append(timeEntries, *tes...)
			return len(*tes), nil
		},
	)

	if err != nil {
		return timeEntries, err
	}

	user, err := c.GetUser(p.UserID)
	if err != nil {
		return timeEntries, err
	}

	for i := range timeEntries {
		timeEntries[i].User = user
	}

	return timeEntries, err
}

func (c *Client) paginate(method, uri string, p PaginationParam, request dto.PaginatedRequest, bodyTempl interface{}, reducer func(interface{}) (int, error)) error {
	page := p.Page
	if p.AllPages {
		page = 1
	}

	if p.PageSize == 0 {
		p.PageSize = 50
	}

	stop := false
	for !stop {
		r, err := c.NewRequest(
			method,
			uri,
			request.WithPagination(page, p.PageSize),
		)
		if err != nil {
			return err
		}

		response := reflect.New(reflect.TypeOf(bodyTempl).Elem()).Interface()
		_, err = c.Do(r, &response)
		if err != nil {
			return err
		}

		count, err := reducer(response)
		if err != nil {
			return err
		}

		stop = count < p.PageSize || !p.AllPages
		page++
	}
	return nil
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

// GetTimeEntryParam params to get a Time Entry
type GetTimeEntryParam struct {
	Workspace   string
	TimeEntryID string
}

// GetTimeEntry will retrieve a Time Entry using its Workspace and ID
func (c *Client) GetTimeEntry(p GetTimeEntryParam) (timeEntry *dto.TimeEntryImpl, err error) {
	r, err := c.NewRequest(
		"GET",
		fmt.Sprintf(
			"v1/workspaces/%s/time-entries/%s",
			p.Workspace,
			p.TimeEntryID,
		),
		nil,
	)

	if err != nil {
		return timeEntry, err
	}

	_, err = c.Do(r, &timeEntry)
	return timeEntry, err
}

// GetTagParam params to find a tag
type GetTagParam struct {
	Workspace string
	TagID     string
}

// GetTag get a single tag, if it exists
func (c *Client) GetTag(p GetTagParam) (*dto.Tag, error) {
	tags, err := c.GetTags(GetTagsParam{
		Workspace: p.Workspace,
	})

	if err != nil {
		return nil, err
	}

	for _, t := range tags {
		if t.ID == p.TagID {
			return &t, nil
		}
	}

	return nil, stackedErrors.Errorf("tag %s not found on workspace %s", p.TagID, p.Workspace)
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

// GetMe get details about the user who created the token
func (c *Client) GetMe() (dto.User, error) {
	r, err := c.NewRequest("GET", "v1/user", nil)

	if err != nil {
		return dto.User{}, err
	}

	var user dto.User
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
	Name      string
	Archived  bool

	PaginationParam
}

// GetTags get all tags of a workspace
func (c *Client) GetTags(p GetTagsParam) ([]dto.Tag, error) {
	var ps, tmpl []dto.Tag

	if p.Workspace == "" {
		return ps, errors.New("workspace needs to be informed to get tags")
	}

	err := c.paginate(
		"GET",
		fmt.Sprintf(
			"v1/workspaces/%s/tags",
			p.Workspace,
		),
		p.PaginationParam,
		dto.GetTagsRequest{
			Name:       p.Name,
			Archived:   p.Archived,
			Pagination: dto.NewPagination(p.Page, p.PageSize),
		},
		&tmpl,
		func(res interface{}) (int, error) {
			if res == nil {
				return 0, nil
			}
			ls := *res.(*[]dto.Tag)

			ps = append(ps, ls...)
			return len(ls), nil
		},
	)
	return ps, err
}

// GetProjectsParam params to get all project of a workspace
type GetProjectsParam struct {
	Workspace string
	Name      string
	Archived  bool

	PaginationParam
}

// GetProjects get all project of a workspace
func (c *Client) GetProjects(p GetProjectsParam) ([]dto.Project, error) {
	var ps, tmpl []dto.Project

	err := c.paginate(
		"GET",
		fmt.Sprintf(
			"v1/workspaces/%s/projects",
			p.Workspace,
		),
		p.PaginationParam,
		dto.GetProjectRequest{
			Name:       p.Name,
			Archived:   p.Archived,
			Pagination: dto.NewPagination(p.Page, p.PageSize),
		},
		&tmpl,
		func(res interface{}) (int, error) {
			if res == nil {
				return 0, nil
			}
			ls := *res.(*[]dto.Project)

			ps = append(ps, ls...)
			return len(ls), nil
		},
	)

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

// DeleteTimeEntryParam params to update a new time entry
type DeleteTimeEntryParam struct {
	Workspace   string
	TimeEntryID string
}

// DeleteTimeEntry deletes a time entry
func (c *Client) DeleteTimeEntry(p DeleteTimeEntryParam) error {
	r, err := c.NewRequest(
		"DELETE",
		fmt.Sprintf(
			"v1/workspaces/%s/time-entries/%s",
			p.Workspace,
			p.TimeEntryID,
		),
		nil,
	)

	if err != nil {
		return err
	}

	_, err = c.Do(r, nil)
	return err
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
