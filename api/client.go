package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/strhlp"
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
	if apiKey == "" {
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

	r, err := c.NewRequest("GET", "v1/workspaces", nil)
	if err != nil {
		return w, err
	}

	_, err = c.Do(r, &w, "GetWorkspaces")

	if err != nil {
		return w, err
	}

	if f.Name == "" {
		return w, nil
	}

	ws := []dto.Workspace{}

	n := strhlp.Normalize(strings.TrimSpace(f.Name))
	for _, i := range w {
		if strings.Contains(strhlp.Normalize(i.Name), n) {
			ws = append(ws, i)
		}
	}

	return ws, nil
}

type field string

const (
	workspaceField   = field("workspace")
	userIDField      = field("user id")
	projectField     = field("project id")
	timeEntryIDField = field("time entry id")
)

func required(action string, values map[field]string) error {
	for f := range values {
		if values[f] == "" {
			return fmt.Errorf("%s is required to %s", f, action)
		}
	}

	return nil
}

type GetWorkspace struct {
	ID string
}

func (c *Client) GetWorkspace(p GetWorkspace) (dto.Workspace, error) {
	err := required("get workspace", map[field]string{
		workspaceField: p.ID,
	})
	if err != nil {
		return dto.Workspace{}, err
	}

	ws, err := c.GetWorkspaces(GetWorkspaces{})
	if err != nil {
		return dto.Workspace{}, err
	}

	for _, w := range ws {
		if w.ID == p.ID {
			return w, nil
		}
	}

	return dto.Workspace{}, dto.Error{Message: "not found", Code: 404}
}

// WorkspaceUsersParam params to query workspace users
type WorkspaceUsersParam struct {
	Workspace string
	Email     string
}

// WorkspaceUsers all users in a Workspace
func (c *Client) WorkspaceUsers(p WorkspaceUsersParam) ([]dto.User, error) {
	var users []dto.User

	err := required("get workspace users", map[field]string{
		workspaceField: p.Workspace,
	})
	if err != nil {
		return users, err
	}

	r, err := c.NewRequest(
		"GET",
		fmt.Sprintf("v1/workspaces/%s/users", p.Workspace),
		nil,
	)
	if err != nil {
		return users, err
	}

	_, err = c.Do(r, &users, "WorkspaceUsers")
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

// AllPages sets the query to retrieve all pages
func AllPages() PaginationParam {
	return PaginationParam{AllPages: true}
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
	Workspace   string
	UserID      string
	FirstDate   time.Time
	LastDate    time.Time
	Description string
	PaginationParam
}

// LogRange list time entries by date range
func (c *Client) LogRange(p LogRangeParam) ([]dto.TimeEntry, error) {
	c.debugf("LogRange - First Date Param: %s | Last Date Param: %s", p.FirstDate, p.LastDate)

	return c.GetUsersHydratedTimeEntries(GetUserTimeEntriesParam{
		Workspace:       p.Workspace,
		UserID:          p.UserID,
		Start:           &p.FirstDate,
		End:             &p.LastDate,
		Description:     p.Description,
		PaginationParam: p.PaginationParam,
	})
}

type GetUserTimeEntriesParam struct {
	Workspace      string
	UserID         string
	OnlyInProgress *bool
	Start          *time.Time
	End            *time.Time
	Description    string

	PaginationParam
}

// GetUserTimeEntries will list the time entries of a user on a workspace, can be paginated
func (c *Client) GetUserTimeEntries(p GetUserTimeEntriesParam) ([]dto.TimeEntryImpl, error) {
	var timeEntries []dto.TimeEntryImpl
	var tes []dto.TimeEntryImpl

	err := c.getUserTimeEntriesImpl(p, false, &tes, func(res interface{}) (int, error) {
		if res == nil {
			return 0, nil
		}

		tes := res.(*[]dto.TimeEntryImpl)
		timeEntries = append(timeEntries, *tes...)
		return len(*tes), nil
	})

	return timeEntries, err
}

// GetUsersHydratedTimeEntries will list hydrated time entries of a user on a workspace, can be paginated
func (c *Client) GetUsersHydratedTimeEntries(p GetUserTimeEntriesParam) ([]dto.TimeEntry, error) {
	var timeEntries []dto.TimeEntry
	var tes []dto.TimeEntry

	err := c.getUserTimeEntriesImpl(p, true, &tes, func(res interface{}) (int, error) {
		if res == nil {
			return 0, nil
		}

		tes := res.(*[]dto.TimeEntry)
		timeEntries = append(timeEntries, *tes...)
		return len(*tes), nil
	})

	if err != nil {
		return timeEntries, err
	}

	user, err := c.GetUser(GetUser{p.Workspace, p.UserID})
	if err != nil {
		return timeEntries, err
	}

	for i := range timeEntries {
		timeEntries[i].User = &user
	}

	return timeEntries, err
}

func (c *Client) getUserTimeEntriesImpl(
	p GetUserTimeEntriesParam,
	hydrated bool,
	tmpl interface{},
	reducer func(interface{}) (int, error),
) error {
	err := required("get time entries", map[field]string{
		workspaceField: p.Workspace,
		userIDField:    p.UserID,
	})
	if err != nil {
		return err
	}
	inProgressFilter := "nil"
	if p.OnlyInProgress != nil {
		if *p.OnlyInProgress {
			inProgressFilter = "true"
		} else {
			inProgressFilter = "false"
		}
	}

	c.debugf("GetUserTimeEntries - Workspace: %s | User: %s | In Progress: %s | Description: %s",
		p.Workspace,
		p.UserID,
		inProgressFilter,
		p.Description,
	)

	r := dto.UserTimeEntriesRequest{
		OnlyInProgress: p.OnlyInProgress,
		Hydrated:       &hydrated,
		Description:    p.Description,
	}

	if p.Start != nil {
		r.Start = &dto.DateTime{Time: *p.Start}
	}

	if p.End != nil {
		r.End = &dto.DateTime{Time: *p.End}
	}

	return c.paginate(
		"GET",
		fmt.Sprintf(
			"v1/workspaces/%s/user/%s/time-entries",
			p.Workspace,
			p.UserID,
		),
		p.PaginationParam,
		r,
		tmpl,
		reducer,
		"GetUserTimeEntries",
	)
}

func (c *Client) paginate(method, uri string, p PaginationParam, request dto.PaginatedRequest, bodyTempl interface{}, reducer func(interface{}) (int, error), name string) error {
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
		_, err = c.Do(r, &response, name)
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

// GetTimeEntryInProgressParam params to query entries
type GetTimeEntryInProgressParam struct {
	Workspace string
	UserID    string
}

// GetTimeEntryInProgress show time entry in progress (if any)
func (c *Client) GetTimeEntryInProgress(p GetTimeEntryInProgressParam) (timeEntryImpl *dto.TimeEntryImpl, err error) {
	b := true
	ts, err := c.GetUserTimeEntries(GetUserTimeEntriesParam{
		Workspace:       p.Workspace,
		UserID:          p.UserID,
		OnlyInProgress:  &b,
		PaginationParam: PaginationParam{PageSize: 1},
	})

	if err != nil {
		return
	}

	if err == nil && len(ts) > 0 {
		timeEntryImpl = &ts[0]
	}
	return
}

// GetHydratedTimeEntryInProgress show hydrated time entry in progress (if any)
func (c *Client) GetHydratedTimeEntryInProgress(p GetTimeEntryInProgressParam) (timeEntry *dto.TimeEntry, err error) {
	b := true
	ts, err := c.GetUsersHydratedTimeEntries(GetUserTimeEntriesParam{
		Workspace:      p.Workspace,
		UserID:         p.UserID,
		OnlyInProgress: &b,
	})
	if err == nil && len(ts) > 0 {
		timeEntry = &ts[0]
	}
	return
}

// GetTimeEntryParam params to get a Time Entry
type GetTimeEntryParam struct {
	Workspace              string
	TimeEntryID            string
	ConsiderDurationFormat bool
}

// GetTimeEntry will retrieve a Time Entry using its Workspace and ID
func (c *Client) GetTimeEntry(p GetTimeEntryParam) (timeEntry *dto.TimeEntryImpl, err error) {
	err = required("get time entry", map[field]string{
		workspaceField:   p.Workspace,
		timeEntryIDField: p.TimeEntryID,
	})

	if err != nil {
		return nil, err
	}

	r, err := c.NewRequest(
		"GET",
		fmt.Sprintf(
			"v1/workspaces/%s/time-entries/%s",
			p.Workspace,
			p.TimeEntryID,
		),
		dto.GetTimeEntryRequest{
			ConsiderDurationFormat: &p.ConsiderDurationFormat,
		},
	)

	if err != nil {
		return timeEntry, err
	}

	_, err = c.Do(r, &timeEntry, "GetTimeEntry")
	return timeEntry, err
}

func (c *Client) GetHydratedTimeEntry(p GetTimeEntryParam) (timeEntry *dto.TimeEntry, err error) {
	err = required("get time entry (hydrated)", map[field]string{
		workspaceField:   p.Workspace,
		timeEntryIDField: p.TimeEntryID,
	})

	if err != nil {
		return nil, err
	}

	b := true
	r, err := c.NewRequest(
		"GET",
		fmt.Sprintf(
			"v1/workspaces/%s/time-entries/%s",
			p.Workspace,
			p.TimeEntryID,
		),
		dto.GetTimeEntryRequest{
			ConsiderDurationFormat: &p.ConsiderDurationFormat,
			Hydrated:               &b,
		},
	)

	if err != nil {
		return timeEntry, err
	}

	_, err = c.Do(r, &timeEntry, "GetHydratedTimeEntry")
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
	err := required("get project", map[field]string{
		workspaceField: p.Workspace,
		projectField:   p.ProjectID,
	})

	if err != nil {
		return project, err
	}

	r, err := c.NewRequest(
		"GET",
		fmt.Sprintf(
			"v1/workspaces/%s/projects/%s",
			p.Workspace,
			p.ProjectID,
		),
		nil,
	)

	if err != nil {
		return project, err
	}

	_, err = c.Do(r, &project, "GetProject")
	return project, err
}

// GetUser params to get a user
type GetUser struct {
	Workspace string
	UserID    string
}

// GetUser filters the wanted user from the workspace users
func (c *Client) GetUser(p GetUser) (dto.User, error) {
	err := required("get user", map[field]string{
		workspaceField: p.Workspace,
		userIDField:    p.UserID,
	})

	if err != nil {
		return dto.User{}, err
	}

	us, err := c.WorkspaceUsers(WorkspaceUsersParam{
		Workspace: p.Workspace,
	})
	if err != nil {
		return dto.User{}, err
	}

	for _, u := range us {
		if u.ID == p.UserID {
			return u, nil
		}
	}

	return dto.User{}, dto.Error{Message: "not found", Code: 404}
}

// GetMe get details about the user who created the token
func (c *Client) GetMe() (dto.User, error) {
	r, err := c.NewRequest("GET", "v1/user", nil)

	if err != nil {
		return dto.User{}, err
	}

	var user dto.User
	_, err = c.Do(r, &user, "GetMe")
	return user, err
}

// GetTasksParam param to find tasks of a project
type GetTasksParam struct {
	Workspace string
	ProjectID string
	Active    bool
	Name      string

	PaginationParam
}

// GetTasks get tasks of a project
func (c *Client) GetTasks(p GetTasksParam) ([]dto.Task, error) {
	var ps, tmpl []dto.Task

	err := required("get tasks", map[field]string{
		workspaceField: p.Workspace,
		projectField:   p.ProjectID,
	})

	if err != nil {
		return ps, err
	}

	err = c.paginate(
		"GET",
		fmt.Sprintf(
			"v1/workspaces/%s/projects/%s/tasks",
			p.Workspace,
			p.ProjectID,
		),
		p.PaginationParam,
		dto.GetTasksRequest{
			Name:   p.Name,
			Active: p.Active,
		},
		&tmpl,
		func(res interface{}) (int, error) {
			if res == nil {
				return 0, nil
			}
			ls := *res.(*[]dto.Task)

			ps = append(ps, ls...)
			return len(ls), nil
		},
		"GetTasks",
	)
	return ps, err
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

	err := required("create time entry", map[field]string{
		workspaceField: p.Workspace,
	})

	if err != nil {
		return t, err
	}

	var end *dto.DateTime
	if p.End != nil {
		end = &dto.DateTime{Time: *p.End}
	}

	r, err := c.NewRequest(
		"POST",
		fmt.Sprintf(
			"v1/workspaces/%s/time-entries",
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

	_, err = c.Do(r, &t, "CreateTimeEntry")
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

	err := required("get tags", map[field]string{
		workspaceField: p.Workspace,
	})

	if err != nil {
		return ps, err
	}

	err = c.paginate(
		"GET",
		fmt.Sprintf(
			"v1/workspaces/%s/tags",
			p.Workspace,
		),
		p.PaginationParam,
		dto.GetTagsRequest{
			Name:     p.Name,
			Archived: p.Archived,
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
		"GetTags",
	)
	return ps, err
}

// GetProjectsParam params to get all project of a workspace
type GetProjectsParam struct {
	Workspace string
	Name      string
	Archived  *bool

	PaginationParam
}

// GetProjects get all project of a workspace
func (c *Client) GetProjects(p GetProjectsParam) ([]dto.Project, error) {
	var ps, tmpl []dto.Project

	err := required("get projects", map[field]string{
		workspaceField: p.Workspace,
	})

	if err != nil {
		return ps, err
	}

	err = c.paginate(
		"GET",
		fmt.Sprintf(
			"v1/workspaces/%s/projects",
			p.Workspace,
		),
		p.PaginationParam,
		dto.GetProjectRequest{
			Name:     p.Name,
			Archived: p.Archived,
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
		"GetProjects",
	)
	return ps, err
}

// OutParam params to end the current time entry
type OutParam struct {
	Workspace string
	UserID    string
	End       time.Time
}

// Out create a new time entry
func (c *Client) Out(p OutParam) error {
	err := required("end running time entry", map[field]string{
		workspaceField: p.Workspace,
		userIDField:    p.UserID,
	})

	if err != nil {
		return err
	}

	r, err := c.NewRequest(
		"PATCH",
		fmt.Sprintf(
			"v1/workspaces/%s/user/%s/time-entries",
			p.Workspace,
			p.UserID,
		),
		dto.OutTimeEntryRequest{
			End: dto.DateTime{Time: p.End},
		},
	)

	if err != nil {
		return err
	}

	_, err = c.Do(r, nil, "Out")
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
	err := required("update time entry", map[field]string{
		workspaceField:   p.Workspace,
		timeEntryIDField: p.TimeEntryID,
	})

	if err != nil {
		return t, err
	}

	var end *dto.DateTime
	if p.End != nil {
		end = &dto.DateTime{Time: *p.End}
	}

	r, err := c.NewRequest(
		"PUT",
		fmt.Sprintf(
			"v1/workspaces/%s/time-entries/%s",
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

	_, err = c.Do(r, &t, "UpdateTimeEntry")
	return t, err
}

// DeleteTimeEntryParam params to update a new time entry
type DeleteTimeEntryParam struct {
	Workspace   string
	TimeEntryID string
}

// DeleteTimeEntry deletes a time entry
func (c *Client) DeleteTimeEntry(p DeleteTimeEntryParam) error {
	err := required("delete time entry", map[field]string{
		workspaceField:   p.Workspace,
		timeEntryIDField: p.TimeEntryID,
	})

	if err != nil {
		return err
	}

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

	_, err = c.Do(r, nil, "DeleteTimeEntry")
	return err
}

type ChangeInvoicedParam struct {
	Workspace    string
	TimeEntryIDs []string
	Invoiced     bool
}

// ChangeInvoiced changes time entries to invoiced or not
func (c *Client) ChangeInvoiced(p ChangeInvoicedParam) error {
	r, err := c.NewRequest(
		"PATCH",
		fmt.Sprintf(
			"v1/workspaces/%s/time-entries/invoiced",
			p.Workspace,
		),
		dto.ChangeTimeEntriesInvoicedRequest{
			TimeEntryIDs: p.TimeEntryIDs,
			Invoiced:     p.Invoiced,
		},
	)

	if err != nil {
		return err
	}

	_, err = c.Do(r, nil, "ChangeInvoiced")
	return err
}
