package api

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/pkg/errors"
)

// Client will help to access Clockify API
type Client interface {
	// SetDebugLogger when set will output the responses of requests to the
	// logger
	SetDebugLogger(logger Logger) Client
	// SetInfoLogger when set will output which requests and params are used to
	// the logger
	SetInfoLogger(logger Logger) Client

	GetWorkspace(GetWorkspace) (dto.Workspace, error)
	GetWorkspaces(GetWorkspaces) ([]dto.Workspace, error)

	GetMe() (dto.User, error)
	GetUser(GetUser) (dto.User, error)
	WorkspaceUsers(WorkspaceUsersParam) ([]dto.User, error)

	AddClient(AddClientParam) (dto.Client, error)
	GetClients(GetClientsParam) ([]dto.Client, error)

	// GetProjects get all project of a workspace
	GetProjects(GetProjectsParam) ([]dto.Project, error)
	// GetProject get a single Project, if exists
	GetProject(GetProjectParam) (*dto.Project, error)
	// AddProject creates a new project
	AddProject(AddProjectParam) (dto.Project, error)
	// UpdateProject changes basic information about the project
	UpdateProject(UpdateProjectParam) (dto.Project, error)
	// UpdateProjectUserCostRate will update the hourly rate of a user on a
	// project
	UpdateProjectUserBillableRate(UpdateProjectUserRateParam) (
		dto.Project, error)
	// UpdateProjectUserCostRate will update the cost of a user on a project
	UpdateProjectUserCostRate(UpdateProjectUserRateParam) (
		dto.Project, error)
	// UpdateProjectEstimate change how the estime of a project is measured
	UpdateProjectEstimate(UpdateProjectEstimateParam) (dto.Project, error)
	// UpdateProjectMemberships changes who has access to add time entries to
	// the project
	UpdateProjectMemberships(UpdateProjectMembershipsParam) (dto.Project, error)
	// UpdateProjectTemplate changes if a project is a template or not
	UpdateProjectTemplate(UpdateProjectTemplateParam) (dto.Project, error)
	// DeleteProject removes a project forever
	DeleteProject(DeleteProjectParam) (dto.Project, error)

	AddTask(AddTaskParam) (dto.Task, error)
	DeleteTask(DeleteTaskParam) (dto.Task, error)
	GetTask(GetTaskParam) (dto.Task, error)
	GetTasks(GetTasksParam) ([]dto.Task, error)
	UpdateTask(UpdateTaskParam) (dto.Task, error)

	GetTag(GetTagParam) (*dto.Tag, error)
	GetTags(GetTagsParam) ([]dto.Tag, error)

	ChangeInvoiced(ChangeInvoicedParam) error
	CreateTimeEntry(CreateTimeEntryParam) (dto.TimeEntryImpl, error)
	DeleteTimeEntry(DeleteTimeEntryParam) error
	GetHydratedTimeEntry(GetTimeEntryParam) (*dto.TimeEntry, error)
	GetHydratedTimeEntryInProgress(GetTimeEntryInProgressParam) (*dto.TimeEntry, error)
	GetTimeEntry(GetTimeEntryParam) (*dto.TimeEntryImpl, error)
	GetTimeEntryInProgress(GetTimeEntryInProgressParam) (*dto.TimeEntryImpl, error)
	GetUserTimeEntries(GetUserTimeEntriesParam) ([]dto.TimeEntryImpl, error)
	GetUsersHydratedTimeEntries(GetUserTimeEntriesParam) ([]dto.TimeEntry, error)
	Log(LogParam) ([]dto.TimeEntry, error)
	LogRange(LogRangeParam) ([]dto.TimeEntry, error)
	UpdateTimeEntry(UpdateTimeEntryParam) (dto.TimeEntryImpl, error)
	Out(OutParam) error
}

type client struct {
	baseURL *url.URL
	http.Client
	debugLogger Logger
	infoLogger  Logger
}

// baseURL is the Clockify API base URL
const baseURL = "https://api.clockify.me/api"

// ErrorMissingAPIKey returned if X-Api-Key is missing
var ErrorMissingAPIKey = errors.New("api Key must be informed")

// ErrorMissingAPIURL returned if base url is missing
var ErrorMissingAPIURL = errors.New("api URL must be informed")

func NewClientFromUrlAndKey(
	apiKey,
	urlString string,
) (Client, error) {
	if apiKey == "" {
		return nil, errors.WithStack(ErrorMissingAPIKey)
	}

	if urlString == "" {
		return nil, errors.WithStack(ErrorMissingAPIURL)
	}

	u, err := url.Parse(urlString)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &client{
		baseURL: u,
		Client: http.Client{
			Transport: transport{
				apiKey: apiKey,
				next:   http.DefaultTransport,
			},
		},
	}, nil
}

// NewClient create a new Client, based on: https://clockify.github.io/clockify_api_docs/
func NewClient(apiKey string) (Client, error) {
	return NewClientFromUrlAndKey(
		apiKey,
		baseURL,
	)
}

// GetWorkspaces will be used to filter the workspaces
type GetWorkspaces struct {
	Name string
}

// Workspaces list all the user's workspaces
func (c *client) GetWorkspaces(f GetWorkspaces) ([]dto.Workspace, error) {
	var w []dto.Workspace

	r, err := c.NewRequest("GET", "v1/workspaces", nil)
	if err != nil {
		return w, err
	}

	_, err = c.Do(r, &w, "GetWorkspaces")

	if err != nil {
		return w, errors.Wrap(err, "get workspaces")
	}

	if f.Name == "" {
		return w, nil
	}

	ws := []dto.Workspace{}

	n := strhlp.Normalize(strings.TrimSpace(f.Name))
	for i := 0; i < len(w); i++ {
		if strings.Contains(strhlp.Normalize(w[i].Name), n) {
			ws = append(ws, w[i])
		}
	}

	return ws, nil
}

type field string

const (
	workspaceField      = field("workspace")
	userIDField         = field("user id")
	userOrGroupIDField  = field("user or group")
	projectField        = field("project id")
	timeEntryIDField    = field("time entry id")
	nameField           = field("name")
	taskIDField         = field("task id")
	estimateMethodField = field("estimate method")
	estimateTypeField   = field("estimate type")
	resetOptionField    = field("reset option")
)

// RequiredFieldError indicates that a field should be filled, but was not
type RequiredFieldError struct {
	Field string
}

func (e RequiredFieldError) Error() string {
	return e.Field + " is required"
}

func required(values map[field]string) error {
	for f := range values {
		if values[f] == "" {
			return RequiredFieldError{Field: string(f)}
		}
	}

	return nil
}

var regexId = regexp.MustCompile("^[a-fA-F0-9]{24}$")

// IsValidID checks if a string looks like a valid ID
func IsValidID(id string) bool {
	return regexId.MatchString(id)
}

// InvalidIDError indicates that a field should be a valid ID, but it is not
type InvalidIDError struct {
	Field string
	ID    string
}

func (e InvalidIDError) Error() string {
	return e.Field + " (\"" + e.ID + "\") is not valid ID"
}

func checkIDs(ids map[field]string) error {
	for field, id := range ids {
		if !IsValidID(id) {
			return InvalidIDError{Field: string(field), ID: id}
		}
	}

	return nil
}

func checkWorkspace(workspace string) error {
	ids := map[field]string{workspaceField: workspace}
	if err := required(ids); err != nil {
		return err
	}

	return checkIDs(ids)
}

func wrapError(err *error, message string, args ...interface{}) {
	if err == nil {
		return
	}
	*err = errors.Wrapf(*err, message, args...)
}

type EntityNotFound struct {
	EntityName string
	ID         string
}

func (e EntityNotFound) Error() string {
	return e.EntityName + " with id " + e.ID + " was not found"
}

func (e EntityNotFound) Unwrap() error {
	return dto.Error{Code: 404, Message: e.Error()}
}

type GetWorkspace struct {
	ID string
}

func (c *client) GetWorkspace(p GetWorkspace) (dto.Workspace, error) {
	var err error
	defer wrapError(&err, "get workspace %s", p.ID)

	if err = checkWorkspace(p.ID); err != nil {
		return dto.Workspace{}, errors.WithStack(err)
	}

	ws, err := c.GetWorkspaces(GetWorkspaces{})
	if err != nil {
		return dto.Workspace{}, err
	}

	for i := 0; i < len(ws); i++ {
		if ws[i].ID == p.ID {
			return ws[i], nil
		}
	}

	err = EntityNotFound{
		EntityName: "workspace",
		ID:         p.ID,
	}

	return dto.Workspace{}, err
}

// WorkspaceUsersParam params to query workspace users
type WorkspaceUsersParam struct {
	Workspace string
	Email     string

	PaginationParam
}

// WorkspaceUsers all users in a Workspace
func (c *client) WorkspaceUsers(p WorkspaceUsersParam) (users []dto.User, err error) {
	defer wrapError(&err, "get users")

	if err := checkWorkspace(p.Workspace); err != nil {
		return users, err
	}

	err = c.paginate(
		"GET",
		fmt.Sprintf("v1/workspaces/%s/users", p.Workspace),
		p.PaginationParam,
		dto.WorkspaceUsersRequest{
			Email: p.Email,
		},
		&users,
		func(res interface{}) (int, error) {
			if res == nil {
				return 0, nil
			}
			ls := *res.(*[]dto.User)

			users = append(users, ls...)
			return len(ls), nil
		},
		"WorkspaceUsers",
	)

	return users, err
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
func (c *client) Log(p LogParam) ([]dto.TimeEntry, error) {
	c.infof("Log - Date Param: %s", p.Date)

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
	ProjectID   string
	TagIDs      []string
	PaginationParam
}

// LogRange list time entries by date range
func (c *client) LogRange(p LogRangeParam) ([]dto.TimeEntry, error) {
	c.infof("LogRange - First Date Param: %s | Last Date Param: %s", p.FirstDate, p.LastDate)

	return c.GetUsersHydratedTimeEntries(GetUserTimeEntriesParam{
		Workspace:       p.Workspace,
		UserID:          p.UserID,
		Start:           &p.FirstDate,
		End:             &p.LastDate,
		Description:     p.Description,
		ProjectID:       p.ProjectID,
		TagIDs:          p.TagIDs,
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
	ProjectID      string
	TagIDs         []string

	PaginationParam
}

// GetUserTimeEntries will list the time entries of a user on a workspace, can be paginated
func (c *client) GetUserTimeEntries(p GetUserTimeEntriesParam) ([]dto.TimeEntryImpl, error) {
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
func (c *client) GetUsersHydratedTimeEntries(p GetUserTimeEntriesParam) ([]dto.TimeEntry, error) {
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

	for i := 0; i < len(timeEntries); i++ {
		timeEntries[i].User = &user
	}

	return timeEntries, err
}

func (c *client) getUserTimeEntriesImpl(
	p GetUserTimeEntriesParam,
	hydrated bool,
	tmpl interface{},
	reducer func(interface{}) (int, error),
) (err error) {
	defer wrapError(&err, "get time entries from user \"%s\"", p.UserID)

	ids := map[field]string{
		workspaceField: p.Workspace,
		userIDField:    p.UserID,
	}

	if err := required(ids); err != nil {
		return err
	}

	if err := checkIDs(ids); err != nil {
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

	c.infof(
		"GetUserTimeEntries - Workspace: %s | User: %s | In Progress: %s "+
			"| Description: %s | Project: %s",
		p.Workspace,
		p.UserID,
		inProgressFilter,
		p.Description,
		p.ProjectID,
	)

	r := dto.UserTimeEntriesRequest{
		OnlyInProgress: p.OnlyInProgress,
		Hydrated:       &hydrated,
		Description:    p.Description,
		Project:        p.ProjectID,
		TagIDs:         p.TagIDs,
	}

	if p.Start != nil {
		r.Start = &dto.DateTime{Time: *p.Start}
	}

	if p.End != nil {
		r.End = &dto.DateTime{Time: *p.End}
	}

	err = c.paginate(
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

	return err
}

func (c *client) paginate(
	method, uri string,
	p PaginationParam,
	request dto.PaginatedRequest,
	bodyTempl interface{},
	reducer func(interface{}) (int, error),
	name string,
) error {
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
func (c *client) GetTimeEntryInProgress(p GetTimeEntryInProgressParam) (timeEntryImpl *dto.TimeEntryImpl, err error) {
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
func (c *client) GetHydratedTimeEntryInProgress(p GetTimeEntryInProgressParam) (timeEntry *dto.TimeEntry, err error) {
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
func (c *client) GetTimeEntry(p GetTimeEntryParam) (timeEntry *dto.TimeEntryImpl, err error) {
	defer wrapError(&err, "get time entry \"%s\"", p.TimeEntryID)

	ids := map[field]string{
		workspaceField:   p.Workspace,
		timeEntryIDField: p.TimeEntryID,
	}

	if err = required(ids); err != nil {
		return nil, err
	}

	if err = checkIDs(ids); err != nil {
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

func (c *client) GetHydratedTimeEntry(p GetTimeEntryParam) (timeEntry *dto.TimeEntry, err error) {
	defer wrapError(&err, "get hydrated time entry \"%s\"", p.TimeEntryID)

	ids := map[field]string{
		workspaceField:   p.Workspace,
		timeEntryIDField: p.TimeEntryID,
	}

	if err = required(ids); err != nil {
		return nil, err
	}

	if err = checkIDs(ids); err != nil {
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
func (c *client) GetTag(p GetTagParam) (*dto.Tag, error) {
	tags, err := c.GetTags(GetTagsParam{
		Workspace: p.Workspace,
	})

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(tags); i++ {
		if tags[i].ID == p.TagID {
			return &tags[i], nil
		}
	}

	return nil, errors.Errorf(
		"tag %s not found on workspace %s", p.TagID, p.Workspace)
}

// GetProjectParam params to get a Project
type GetProjectParam struct {
	Workspace string
	ProjectID string
	Hydrate   bool
}

// GetProject get a single Project, if exists
func (c *client) GetProject(p GetProjectParam) (pr *dto.Project, err error) {
	defer wrapError(&err, "get project \"%s\"", p.ProjectID)

	ids := map[field]string{
		workspaceField: p.Workspace,
		projectField:   p.ProjectID,
	}

	if err = required(ids); err != nil {
		return pr, err
	}

	if err = checkIDs(ids); err != nil {
		return pr, err
	}

	r, err := c.NewRequest(
		"GET",
		fmt.Sprintf(
			"v1/workspaces/%s/projects/%s",
			p.Workspace,
			p.ProjectID,
		),
		dto.GetProjectRequest{Hydrated: p.Hydrate},
	)

	if err != nil {
		return pr, err
	}

	_, err = c.Do(r, &pr, "GetProject")
	if p.Hydrate && pr != nil {
		pr.Hydrated = true
	}

	return pr, err
}

// GetUser params to get a user
type GetUser struct {
	Workspace string
	UserID    string
}

// GetUser filters the wanted user from the workspace users
func (c *client) GetUser(p GetUser) (dto.User, error) {
	var err error
	defer wrapError(&err, "get user \"%s\"", p.UserID)

	ids := map[field]string{
		workspaceField: p.Workspace,
		userIDField:    p.UserID,
	}

	if err = required(ids); err != nil {
		return dto.User{}, err
	}

	if err = checkIDs(ids); err != nil {
		return dto.User{}, err
	}

	us, err := c.WorkspaceUsers(WorkspaceUsersParam{
		Workspace:       p.Workspace,
		PaginationParam: AllPages(),
	})
	if err != nil {
		return dto.User{}, errors.Wrapf(err, "get user %s", p.UserID)
	}

	for i := 0; i < len(us); i++ {
		if us[i].ID == p.UserID {
			return us[i], nil
		}
	}

	err = EntityNotFound{
		EntityName: "user",
		ID:         p.UserID,
	}
	return dto.User{}, err
}

// GetMe get details about the user who created the token
func (c *client) GetMe() (dto.User, error) {
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
func (c *client) GetTasks(p GetTasksParam) (ps []dto.Task, err error) {
	var tmpl []dto.Task

	defer wrapError(&err, "get tasks from project \"%s\"", p.ProjectID)

	ids := map[field]string{
		workspaceField: p.Workspace,
		projectField:   p.ProjectID,
	}

	if err = required(ids); err != nil {
		return ps, err
	}

	if err = checkIDs(ids); err != nil {
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

// GetTaskParam param to get a task on a project
type GetTaskParam struct {
	Workspace string
	ProjectID string
	TaskID    string
}

// GetTasks get tasks of a project
func (c *client) GetTask(p GetTaskParam) (t dto.Task, err error) {
	defer wrapError(&err, "get task \"%s\"", p.TaskID)

	ids := map[field]string{
		workspaceField: p.Workspace,
		projectField:   p.ProjectID,
		taskIDField:    p.TaskID,
	}

	if err = required(ids); err != nil {
		return t, err
	}

	if err = checkIDs(ids); err != nil {
		return t, err
	}

	r, err := c.NewRequest(
		"GET",
		fmt.Sprintf(
			"v1/workspaces/%s/projects/%s/tasks/%s",
			p.Workspace,
			p.ProjectID,
			p.TaskID,
		),
		nil,
	)

	if err != nil {
		return t, err
	}

	_, err = c.Do(r, &t, "GetTask")
	return t, err
}

type TaskStatus string

const (
	TaskStatusDefault = ""
	TaskStatusDone    = "DONE"
	TaskStatusActive  = "ACTIVE"
)

// AddTaskParam param to add tasks to a project
type AddTaskParam struct {
	Workspace   string
	ProjectID   string
	Name        string
	AssigneeIDs *[]string
	Estimate    *time.Duration
	Status      TaskStatus
	Billable    *bool
}

func (c *client) AddTask(p AddTaskParam) (task dto.Task, err error) {
	defer wrapError(&err, "add task to project \"%s\"", p.ProjectID)

	if err = required(map[field]string{
		nameField:      p.Name,
		workspaceField: p.Workspace,
		projectField:   p.ProjectID,
	}); err != nil {
		return task, err
	}

	if err = checkIDs(map[field]string{
		workspaceField: p.Workspace,
		projectField:   p.ProjectID,
	}); err != nil {
		return task, err
	}

	r := dto.AddTaskRequest{
		Name:        p.Name,
		AssigneeIDs: p.AssigneeIDs,
		Billable:    p.Billable,
	}

	if p.Status != TaskStatus("") {
		s := string(p.Status)
		r.Status = &s
	}

	if p.Estimate != nil {
		e := dto.Duration{Duration: *p.Estimate}
		r.Estimate = &e
	}

	req, err := c.NewRequest(
		"POST",
		fmt.Sprintf(
			"v1/workspaces/%s/projects/%s/tasks",
			p.Workspace,
			p.ProjectID,
		),
		r,
	)

	if err != nil {
		return task, err
	}

	_, err = c.Do(req, &task, "AddTask")
	return task, err
}

// UpdateTaskParam param to update tasks to a project
type UpdateTaskParam struct {
	Workspace   string
	ProjectID   string
	TaskID      string
	Name        string
	AssigneeIDs *[]string
	Estimate    *time.Duration
	Status      TaskStatus
	Billable    *bool
}

func (c *client) UpdateTask(p UpdateTaskParam) (task dto.Task, err error) {
	defer wrapError(&err, "update task \"%s\"", p.TaskID)

	if err = required(map[field]string{
		nameField:      p.Name,
		taskIDField:    p.TaskID,
		workspaceField: p.Workspace,
		projectField:   p.ProjectID,
	}); err != nil {
		return task, err
	}

	if err = checkIDs(map[field]string{
		taskIDField:    p.TaskID,
		workspaceField: p.Workspace,
		projectField:   p.ProjectID,
	}); err != nil {
		return task, err
	}

	r := dto.UpdateTaskRequest{
		Name:        p.Name,
		AssigneeIDs: p.AssigneeIDs,
		Billable:    p.Billable,
	}

	if p.Status != TaskStatus("") {
		s := string(p.Status)
		r.Status = &s
	}

	if p.Estimate != nil {
		e := dto.Duration{Duration: *p.Estimate}
		r.Estimate = &e
	}

	req, err := c.NewRequest(
		"PUT",
		fmt.Sprintf(
			"v1/workspaces/%s/projects/%s/tasks/%s",
			p.Workspace,
			p.ProjectID,
			p.TaskID,
		),
		r,
	)

	if err != nil {
		return task, err
	}

	_, err = c.Do(req, &task, "UpdateTask")
	return task, err
}

// DeleteTaskParam param to update tasks to a project
type DeleteTaskParam struct {
	Workspace string
	ProjectID string
	TaskID    string
}

func (c *client) DeleteTask(p DeleteTaskParam) (task dto.Task, err error) {
	defer wrapError(&err, "delete task \"%s\"", p.TaskID)

	ids := map[field]string{
		taskIDField:    p.TaskID,
		workspaceField: p.Workspace,
		projectField:   p.ProjectID,
	}

	if err = required(ids); err != nil {
		return task, err
	}

	if err = checkIDs(ids); err != nil {
		return task, err
	}

	req, err := c.NewRequest(
		"DELETE",
		fmt.Sprintf(
			"v1/workspaces/%s/projects/%s/tasks/%s",
			p.Workspace,
			p.ProjectID,
			p.TaskID,
		),
		nil,
	)

	if err != nil {
		return task, err
	}

	_, err = c.Do(req, &task, "DeleteTask")
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
func (c *client) CreateTimeEntry(p CreateTimeEntryParam) (
	t dto.TimeEntryImpl, err error) {
	defer wrapError(&err, "create time entry")

	if err = checkWorkspace(p.Workspace); err != nil {
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
	Archived  *bool

	PaginationParam
}

// GetTags get all tags of a workspace
func (c *client) GetTags(p GetTagsParam) (ps []dto.Tag, err error) {
	defer wrapError(&err, "get tags")
	var tmpl []dto.Tag
	if err = checkWorkspace(p.Workspace); err != nil {
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

// GetClientsParam params to get all clients of a workspace
type GetClientsParam struct {
	Workspace string
	Name      string
	Archived  *bool

	PaginationParam
}

// GetClients gets all clients of a workspace
func (c *client) GetClients(p GetClientsParam) (
	clients []dto.Client, err error) {
	defer wrapError(&err, "get clients")

	var tmpl []dto.Client
	if err = checkWorkspace(p.Workspace); err != nil {
		return clients, err
	}

	err = c.paginate(
		"GET",
		fmt.Sprintf(
			"v1/workspaces/%s/clients",
			p.Workspace,
		),
		p.PaginationParam,
		dto.GetClientsRequest{
			Name:     p.Name,
			Archived: p.Archived,
		},
		&tmpl,
		func(res interface{}) (int, error) {
			if res == nil {
				return 0, nil
			}
			ls := *res.(*[]dto.Client)

			clients = append(clients, ls...)
			return len(ls), nil
		},
		"GetClients",
	)
	return clients, err
}

type AddClientParam struct {
	Workspace string
	Name      string
}

// AddClient adds a new client to a workspace
func (c *client) AddClient(p AddClientParam) (client dto.Client, err error) {
	defer wrapError(&err, "add client")

	if err = required(map[field]string{
		nameField:      p.Name,
		workspaceField: p.Workspace,
	}); err != nil {
		return client, err
	}

	if err = checkIDs(map[field]string{
		workspaceField: p.Workspace,
	}); err != nil {
		return client, err
	}

	req, err := c.NewRequest(
		"POST",
		fmt.Sprintf(
			"v1/workspaces/%s/clients",
			p.Workspace,
		),
		dto.AddClientRequest{
			Name: p.Name,
		},
	)

	if err != nil {
		return client, err
	}

	_, err = c.Do(req, &client, "AddClient")
	return client, err
}

// GetProjectsParam params to get all project of a workspace
type GetProjectsParam struct {
	Workspace string
	Name      string
	Clients   []string
	Archived  *bool
	Hydrate   bool

	PaginationParam
}

// GetProjects get all project of a workspace
func (c *client) GetProjects(p GetProjectsParam) (ps []dto.Project, err error) {
	defer wrapError(&err, "get projects")

	var tmpl []dto.Project
	if err = checkWorkspace(p.Workspace); err != nil {
		return ps, err
	}

	err = c.paginate(
		"GET",
		fmt.Sprintf(
			"v1/workspaces/%s/projects",
			p.Workspace,
		),
		p.PaginationParam,
		dto.GetProjectsRequest{
			Name:     p.Name,
			Archived: p.Archived,
			Clients:  p.Clients,
			Hydrated: p.Hydrate,
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

	if p.Hydrate {
		for i := range ps {
			ps[i].Hydrated = true
		}
	}

	return ps, err
}

type AddProjectParam struct {
	Workspace string
	Name      string
	ClientId  string
	Color     string
	Note      string
	Billable  bool
	Public    bool
}

func parseColor(c string) (string, error) {
	if !strings.HasPrefix(c, "#") {
		c = "#" + c
	}

	if len(c) != 4 && len(c) != 7 {
		return c, errors.New("color must have 3 (#000) or 6 (#ffffff) numbers")
	}

	if len(c) == 4 {
		c = string([]byte{'#', c[1], c[1], c[2], c[2], c[3], c[3]})
	}

	if _, err := hex.DecodeString(c[1:]); err != nil {
		return c, errors.Wrap(err, "color \""+c+"\" is not a hex string")
	}

	return c, nil
}

// AddProject adds a new project to a workspace
func (c *client) AddProject(p AddProjectParam) (
	project dto.Project, err error) {
	defer wrapError(&err, "add project")

	if err = required(map[field]string{
		nameField:      p.Name,
		workspaceField: p.Workspace,
	}); err != nil {
		return project, err
	}

	if err = checkIDs(map[field]string{
		workspaceField: p.Workspace,
	}); err != nil {
		return project, err
	}

	if p.Color != "" {
		p.Color, err = parseColor(p.Color)
		if err != nil {
			return
		}
	}

	req, err := c.NewRequest(
		"POST",
		fmt.Sprintf(
			"v1/workspaces/%s/projects",
			p.Workspace,
		),
		dto.AddProjectRequest{
			Name:     p.Name,
			ClientId: p.ClientId,
			IsPublic: p.Public,
			Color:    p.Color,
			Note:     p.Note,
			Billable: p.Billable,
			Public:   p.Public,
		},
	)

	if err != nil {
		return project, err
	}

	_, err = c.Do(req, &project, "AddProject")
	return project, err
}

// UpdateProjectParam sets the properties to change on a project
// Workspace and ID are required
type UpdateProjectParam struct {
	Workspace string
	ProjectID string
	Name      string
	ClientId  *string
	Color     string
	Note      *string
	Billable  *bool
	Public    *bool
	Archived  *bool
}

// UpdateProject will change properties of a Project, leave the property as nil
// or "empty" to not change it
func (c *client) UpdateProject(p UpdateProjectParam) (
	project dto.Project, err error) {
	defer wrapError(&err, "update project")

	if err = required(map[field]string{
		projectField:   p.ProjectID,
		workspaceField: p.Workspace,
	}); err != nil {
		return project, err
	}

	if err = checkIDs(map[field]string{
		projectField:   p.ProjectID,
		workspaceField: p.Workspace,
	}); err != nil {
		return project, err
	}

	if p.Color != "" {
		p.Color, err = parseColor(p.Color)
		if err != nil {
			return
		}
	}

	var name, color *string
	if p.Name != "" {
		name = &p.Name
	}
	if p.Color != "" {
		color = &p.Color
	}

	req, err := c.NewRequest(
		"PUT",
		"v1/workspaces/"+p.Workspace+"/projects/"+p.ProjectID,
		dto.UpdateProjectRequest{
			Name:     name,
			ClientId: p.ClientId,
			IsPublic: p.Public,
			Color:    color,
			Note:     p.Note,
			Billable: p.Billable,
			Archived: p.Archived,
		},
	)

	if err != nil {
		return project, err
	}

	_, err = c.Do(req, &project, "UpdateProject")
	return project, err
}

// UpdateMembership represents the membership of a User or User Group to a
// project
type UpdateMembership struct {
	UserOrGroupID    string
	HourlyRateAmount int64
}

// UpdateProjectMembershipsParam will change which users and groups have
// access to the project
type UpdateProjectMembershipsParam struct {
	Workspace   string
	ProjectID   string
	Memberships []UpdateMembership
}

// UpdateProjectMemberships changes who has access to add time entries to
// the project
func (c *client) UpdateProjectMemberships(p UpdateProjectMembershipsParam) (
	pr dto.Project, err error) {
	defer wrapError(&err, "update project memberships")

	if err = required(map[field]string{
		projectField:   p.ProjectID,
		workspaceField: p.Workspace,
	}); err != nil {
		return
	}

	if err = checkIDs(map[field]string{
		projectField:   p.ProjectID,
		workspaceField: p.Workspace,
	}); err != nil {
		return
	}

	members := make([]dto.UpdateProjectMembership, len(p.Memberships))
	for i := range p.Memberships {
		id := map[field]string{
			userOrGroupIDField: p.Memberships[i].UserOrGroupID}
		if err = required(id); err != nil {
			return
		}

		if err = checkIDs(id); err != nil {
			return
		}

		members[i].UserID = p.Memberships[i].UserOrGroupID
		members[i].HourlyRate.Amount = p.Memberships[i].HourlyRateAmount
	}

	req, err := c.NewRequest(
		"PATCH",
		"v1/workspaces/"+p.Workspace+"/projects/"+p.ProjectID+"/memberships",
		dto.UpdateProjectMembershipsRequest{
			Memberships: members,
		},
	)

	if err != nil {
		return pr, err
	}

	_, err = c.Do(req, &pr, "UpdateProjectMemberships")
	return pr, err
}

// UpdateProjectTemplateParam sets which project will be updated,and if it will
// became a template or not
type UpdateProjectTemplateParam struct {
	Workspace string
	ProjectID string
	Template  bool
}

// UpdateProjectTemplate changes if a project is a template or not
func (c *client) UpdateProjectTemplate(p UpdateProjectTemplateParam) (
	pr dto.Project, err error) {
	defer wrapError(&err, "update project template")

	if err = required(map[field]string{
		projectField:   p.ProjectID,
		workspaceField: p.Workspace,
	}); err != nil {
		return
	}

	if err = checkIDs(map[field]string{
		projectField:   p.ProjectID,
		workspaceField: p.Workspace,
	}); err != nil {
		return
	}

	req, err := c.NewRequest(
		"PATCH",
		"v1/workspaces/"+p.Workspace+"/projects/"+p.ProjectID+"/template",
		dto.UpdateProjectTemplateRequest{
			IsTemplate: p.Template,
		},
	)

	if err != nil {
		return pr, err
	}

	_, err = c.Do(req, &pr, "UpdateProjectTemplate")
	return pr, err
}

// UpdateProjectUserRateParam sets the parameters to update the billable/cost
// rate, if Since is not nil, then all time entries after that time will be
// updated to new rate
type UpdateProjectUserRateParam struct {
	Workspace string
	ProjectID string
	UserID    string
	Amount    uint
	Since     *time.Time
}

func (c *client) UpdateProjectUserBillableRate(
	p UpdateProjectUserRateParam) (project dto.Project, err error) {
	defer wrapError(&err, "update project user billable rate")

	if err = required(map[field]string{
		projectField:   p.ProjectID,
		workspaceField: p.Workspace,
		userIDField:    p.UserID,
	}); err != nil {
		return
	}

	if err = checkIDs(map[field]string{
		projectField:   p.ProjectID,
		workspaceField: p.Workspace,
		userIDField:    p.UserID,
	}); err != nil {
		return
	}

	var since *dto.DateTime
	if p.Since != nil {
		since = &dto.DateTime{Time: *p.Since}
	}

	req, err := c.NewRequest(
		"PUT",
		"v1/workspaces/"+p.Workspace+"/projects/"+p.ProjectID+
			"/users/"+p.UserID+"/hourly-rate",
		dto.UpdateProjectUserRateRequest{
			Amount: p.Amount,
			Since:  since,
		},
	)

	if err != nil {
		return project, err
	}

	_, err = c.Do(req, &project, "UpdateProjectUserBillableRate")
	return project, err
}

func (c *client) UpdateProjectUserCostRate(
	p UpdateProjectUserRateParam) (project dto.Project, err error) {
	defer wrapError(&err, "update project user cost rate")

	if err = required(map[field]string{
		projectField:   p.ProjectID,
		workspaceField: p.Workspace,
		userIDField:    p.UserID,
	}); err != nil {
		return
	}

	if err = checkIDs(map[field]string{
		projectField:   p.ProjectID,
		workspaceField: p.Workspace,
		userIDField:    p.UserID,
	}); err != nil {
		return
	}

	var since *dto.DateTime
	if p.Since != nil {
		since = &dto.DateTime{Time: *p.Since}
	}

	req, err := c.NewRequest(
		"PUT",
		"v1/workspaces/"+p.Workspace+"/projects/"+p.ProjectID+
			"/users/"+p.UserID+"/cost-rate",
		dto.UpdateProjectUserRateRequest{
			Amount: p.Amount,
			Since:  since,
		},
	)

	if err != nil {
		return project, err
	}

	_, err = c.Do(req, &project, "UpdateProjectUserCostRate")
	return project, err
}

// EstimateMethod are methods to estimate projects (none, budget and time)
type EstimateMethod string

const (
	// EstimateMethodNone dont estimate the project
	EstimateMethodNone = EstimateMethod("none")
	// EstimateMethodTime estimate by time
	EstimateMethodTime = EstimateMethod("time")
	// EstimateMethodBudget estimate by budget
	EstimateMethodBudget = EstimateMethod("budget")
)

// EstimateType sets if the estimate is for the role project or per task
type EstimateType string

const (
	EstimateTypeProject = EstimateType("project")
	EstimateTypeTask    = EstimateType("task")
)

func (t EstimateType) toRequestType() *dto.EstimateType {
	switch t {
	case EstimateTypeTask:
		v := dto.EstimateTypeAuto
		return &v
	case EstimateTypeProject:
		v := dto.EstimateTypeManual
		return &v
	default:
		return nil
	}
}

// EstimateResetOption defines the period in which the estimates reset
type EstimateResetOption string

const (
	EstimateResetOptionDefault = EstimateType("")
	EstimateResetOptionMonthly = EstimateResetOption("monthly")
)

func (t EstimateResetOption) toRequestType() *dto.EstimateResetOption {
	switch t {
	case EstimateResetOptionMonthly:
		v := dto.EstimateResetOptionMonthly
		return &v
	default:
		return nil
	}
}

// UpdateProjectEstimateParam holds parameters to change project estimate
type UpdateProjectEstimateParam struct {
	Workspace   string
	ProjectID   string
	Method      EstimateMethod
	Type        EstimateType
	ResetOption EstimateResetOption
	Estimate    int64
}

// UpdateProjectEstimate change how the estime of a project is measured
func (c *client) UpdateProjectEstimate(p UpdateProjectEstimateParam) (
	r dto.Project, err error) {
	defer wrapError(&err, "update project estimate")

	if err = required(map[field]string{
		projectField:        p.ProjectID,
		workspaceField:      p.Workspace,
		estimateMethodField: string(p.Method),
	}); err != nil {
		return
	}

	if err = checkIDs(map[field]string{
		projectField:   p.ProjectID,
		workspaceField: p.Workspace,
	}); err != nil {
		return
	}

	if err = shouldBeOneOf(estimateMethodField, string(p.Method), []string{
		string(EstimateMethodNone),
		string(EstimateMethodTime),
		string(EstimateMethodBudget),
	}); err != nil {
		return
	}

	if p.Method != EstimateMethodNone {
		if err = shouldBeOneOf(estimateTypeField, string(p.Type), []string{
			string(EstimateTypeProject),
			string(EstimateTypeTask),
		}); err != nil {
			return
		}

		if err = shouldBeOneOf(resetOptionField, string(p.ResetOption),
			[]string{
				string(EstimateResetOptionDefault),
				string(EstimateResetOptionMonthly),
			}); err != nil {
			return
		}

		if p.Type != EstimateTypeProject {
			p.Estimate = 0
		} else if p.Estimate <= 0 {
			err = errors.New(
				"estimate should be greater than zero for type project")
			return
		}
	}

	b := dto.UpdateProjectEstimateRequest{}
	if p.Method != EstimateMethodNone {
		be := dto.BaseEstimateRequest{
			Active:       true,
			Type:         p.Type.toRequestType(),
			ResetOptions: p.ResetOption.toRequestType(),
		}

		switch p.Method {
		case EstimateMethodBudget:
			b.BudgetEstimate.BaseEstimateRequest = be
			if p.Estimate > 0 {
				e := uint64(p.Estimate)
				b.BudgetEstimate.Estimate = &e
			}
		case EstimateMethodTime:
			b.TimeEstimate.BaseEstimateRequest = be
			if p.Estimate > 0 {
				b.TimeEstimate.Estimate = &dto.Duration{
					Duration: time.Duration(p.Estimate)}
			}
		}
	}

	req, err := c.NewRequest(
		"PATCH",
		"v1/workspaces/"+p.Workspace+"/projects/"+p.ProjectID+"/estimate",
		b,
	)

	if err != nil {
		return
	}

	_, err = c.Do(req, &r, "UpdateProjectEstimate")

	return
}

// DeleteProjectParam identifies which project to delete
type DeleteProjectParam struct {
	Workspace string
	ProjectID string
}

// DeleteProject removes a project forever
func (c *client) DeleteProject(p DeleteProjectParam) (
	pr dto.Project, err error) {
	defer wrapError(&err, "delete project")

	ids := map[field]string{
		workspaceField: p.Workspace,
		projectField:   p.ProjectID,
	}

	if err = required(ids); err != nil {
		return pr, err
	}

	if err = checkIDs(ids); err != nil {
		return pr, err
	}

	r, err := c.NewRequest(
		"DELETE",
		"v1/workspaces/"+p.Workspace+"/projects/"+p.ProjectID,
		nil,
	)

	if err != nil {
		return pr, err
	}

	_, err = c.Do(r, &pr, "DeleteProject")

	return pr, err
}

// InvalidOptionError indicates that the parameter has a limited set of valid
// values, and the one used is not one of them (see Options for the valid ones)
type InvalidOptionError struct {
	Field   string
	Options []string
}

func (i *InvalidOptionError) Error() string {
	return "valid options for " + i.Field + " are " + strhlp.ListForHumans(i.Options)
}

func shouldBeOneOf(f field, s string, o []string) error {
	if strhlp.InSlice(s, o) {
		return nil
	}

	return &InvalidOptionError{
		Field:   string(f),
		Options: o,
	}
}

// OutParam params to end the current time entry
type OutParam struct {
	Workspace string
	UserID    string
	End       time.Time
}

// Out create a new time entry
func (c *client) Out(p OutParam) (err error) {
	defer wrapError(&err, "end running time entry")

	ids := map[field]string{
		workspaceField: p.Workspace,
		userIDField:    p.UserID,
	}

	if err = required(ids); err != nil {
		return err
	}

	if err = checkIDs(ids); err != nil {
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
func (c *client) UpdateTimeEntry(p UpdateTimeEntryParam) (
	t dto.TimeEntryImpl, err error) {
	defer wrapError(&err, "update time entry \"%s\"", p.TimeEntryID)

	ids := map[field]string{
		workspaceField:   p.Workspace,
		timeEntryIDField: p.TimeEntryID,
	}

	if err = required(ids); err != nil {
		return t, err
	}

	if err = checkIDs(ids); err != nil {
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
func (c *client) DeleteTimeEntry(p DeleteTimeEntryParam) (err error) {
	defer wrapError(&err, "delete time entry \"%s\"", p.TimeEntryID)

	ids := map[field]string{
		workspaceField:   p.Workspace,
		timeEntryIDField: p.TimeEntryID,
	}

	if err = required(ids); err != nil {
		return err
	}

	if err = checkIDs(ids); err != nil {
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
func (c *client) ChangeInvoiced(p ChangeInvoicedParam) error {
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
