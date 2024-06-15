package dto

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// DateTime is a time presentation for parameters
type DateTime struct {
	time.Time
}

// MarshalJSON converts DateTime correctly
func (d DateTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(d.String())), nil
}

func (d DateTime) String() string {
	return d.Time.UTC().Format("2006-01-02T15:04:05Z")
}

// Duration is a time presentation for parameters
type Duration struct {
	time.Duration
}

// MarshalJSON converts Duration correctly
func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.String() + "\""), nil
}

// UnmarshalJSON converts a JSON value to Duration correctly
func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return errors.Wrap(err, "unmarshal duration")
	}

	dc, err := StringToDuration(s)
	if err != nil {
		return err
	}

	*d = Duration{dc}
	return err
}

func StringToDuration(s string) (time.Duration, error) {
	if len(s) < 4 {
		return 0, errors.Errorf("duration %s is invalid", s)
	}

	var u, dc time.Duration
	var j, i int
	for ; i < len(s); i++ {
		switch s[i] {
		case 'P', 'T':
			j = i + 1
			continue
		case 'H':
			u = time.Hour
		case 'M':
			u = time.Minute
		case 'S':
			u = time.Second
		default:
			continue
		}

		v, err := strconv.Atoi(s[j:i])
		if err != nil {
			return 0, errors.Wrap(err, "cast cast "+s[j:i]+" to int")
		}
		dc = dc + time.Duration(v)*u
		j = i + 1
	}

	return dc, nil
}

func (d Duration) String() string {
	s := d.Duration.String()
	i := strings.LastIndex(s, ".")
	if i > -1 {
		s = s[0:i] + "s"
	}

	return "PT" + strings.ToUpper(s)
}

func (dd Duration) HumanString() string {
	d := dd.Duration
	p := ""
	if d < 0 {
		p = "-"
		d = d * -1
	}

	return p + fmt.Sprintf("%d:%02d:%02d",
		int64(d.Hours()), int64(d.Minutes())%60, int64(d.Seconds())%60)
}

type pagination struct {
	page     int
	pageSize int
}

func newPagination(page, size int) pagination {
	return pagination{
		page:     page,
		pageSize: size,
	}
}

// AppendToQuery decorates the URL with pagination parameters
func (p pagination) AppendToQuery(u *url.URL) *url.URL {
	v := u.Query()

	if p.page != 0 {
		v.Add("page", strconv.Itoa(p.page))
	}
	if p.pageSize != 0 {
		v.Add("page-size", strconv.Itoa(p.pageSize))
	}

	u.RawQuery = v.Encode()

	return u
}

type PaginatedRequest interface {
	WithPagination(page, size int) PaginatedRequest
}

// GetTimeEntryRequest to get a time entry
type GetTimeEntryRequest struct {
	Hydrated               *bool
	ConsiderDurationFormat *bool
}

// AppendToQuery decorates the URL with the query string needed for this Request
func (r GetTimeEntryRequest) AppendToQuery(u *url.URL) *url.URL {
	v := u.Query()
	if r.Hydrated != nil && *r.Hydrated {
		v.Add("hydrated", "true")
	}
	if r.ConsiderDurationFormat != nil && *r.ConsiderDurationFormat {
		v.Add("consider-duration-format", "true")
	}

	u.RawQuery = v.Encode()

	return u
}

// UserTimeEntriesRequest to get entries of a user
type UserTimeEntriesRequest struct {
	Description string
	Start       *DateTime
	End         *DateTime
	Project     string
	Task        string
	TagIDs      []string

	ProjectRequired        *bool
	TaskRequired           *bool
	ConsiderDurationFormat *bool
	Hydrated               *bool
	OnlyInProgress         *bool

	pagination
}

// WithPagination add pagination to the UserTimeEntriesRequest
func (r UserTimeEntriesRequest) WithPagination(page, size int) PaginatedRequest {
	r.pagination = newPagination(page, size)
	return r
}

// AppendToQuery decorates the URL with the query string needed for this Request
func (r UserTimeEntriesRequest) AppendToQuery(u *url.URL) *url.URL {
	u = r.pagination.AppendToQuery(u)
	v := u.Query()

	if r.Start != nil {
		v.Add("start", r.Start.String())
	}

	if r.End != nil {
		v.Add("end", r.End.String())
	}

	addNotNil := func(b *bool, p string) {
		if b == nil {
			return
		}

		if *b {
			v.Add(p, "1")
		} else {
			v.Add(p, "0")
		}

	}

	addNotNil(r.ProjectRequired, "project-required")
	addNotNil(r.TaskRequired, "task-required")
	addNotNil(r.ConsiderDurationFormat, "consider-duration-format")
	addNotNil(r.Hydrated, "hydrated")
	addNotNil(r.OnlyInProgress, "in-progress")

	addNotEmpty := func(s string, p string) {
		if s == "" {
			return
		}

		v.Add(p, s)
	}

	addNotEmpty(r.Description, "description")
	addNotEmpty(r.Project, "project")
	addNotEmpty(r.Task, "task")

	for _, t := range r.TagIDs {
		addNotEmpty(t, "tags")
	}

	u.RawQuery = v.Encode()

	return u
}

// OutTimeEntryRequest to end the current time entry
type OutTimeEntryRequest struct {
	End DateTime `json:"end"`
}

// CreateTimeEntryRequest to create a time entry is created
type CreateTimeEntryRequest struct {
	Start        DateTime           `json:"start,omitempty"`
	End          *DateTime          `json:"end,omitempty"`
	Billable     *bool              `json:"billable,omitempty"`
	Description  string             `json:"description,omitempty"`
	ProjectID    string             `json:"projectId,omitempty"`
	TaskID       string             `json:"taskId,omitempty"`
	TagIDs       []string           `json:"tagIds,omitempty"`
	CustomFields []CustomFieldValue `json:"customFields,omitempty"`
}

// CustomFieldValue DTO
type CustomFieldValue struct {
	CustomFieldID string `json:"customFieldId"`
	Status        string `json:"status"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Value         string `json:"value"`
}

// UpdateTimeEntryRequest to update a time entry
type UpdateTimeEntryRequest struct {
	Start        DateTime           `json:"start,omitempty"`
	End          *DateTime          `json:"end,omitempty"`
	Billable     bool               `json:"billable,omitempty"`
	Description  string             `json:"description,omitempty"`
	ProjectID    string             `json:"projectId,omitempty"`
	TaskID       string             `json:"taskId,omitempty"`
	TagIDs       []string           `json:"tagIds,omitempty"`
	CustomFields []CustomFieldValue `json:"customFields,omitempty"`
}

type GetClientsRequest struct {
	Name     string
	Archived *bool

	pagination
}

// WithPagination add pagination to the GetClientsRequest
func (r GetClientsRequest) WithPagination(page, size int) PaginatedRequest {
	r.pagination = newPagination(page, size)
	return r
}

// AppendToQuery decorates the URL with the query string needed for this Request
func (r GetClientsRequest) AppendToQuery(u *url.URL) *url.URL {
	u = r.pagination.AppendToQuery(u)

	v := u.Query()

	if r.Name != "" {
		v.Add("name", r.Name)
	}

	if r.Archived != nil {
		v.Add("archived", boolString[*r.Archived])
	}

	u.RawQuery = v.Encode()

	return u
}

type AddClientRequest struct {
	Name string `json:"name"`
}

type GetProjectsRequest struct {
	Name     string
	Archived *bool
	Clients  []string
	Hydrated bool

	pagination
}

// WithPagination add pagination to the GetProjectRequest
func (r GetProjectsRequest) WithPagination(page, size int) PaginatedRequest {
	r.pagination = newPagination(page, size)
	return r
}

var boolString = map[bool]string{
	true:  "true",
	false: "false",
}

// AppendToQuery decorates the URL with the query string needed for this Request
func (r GetProjectsRequest) AppendToQuery(u *url.URL) *url.URL {
	u = r.pagination.AppendToQuery(u)

	v := u.Query()

	if r.Name != "" {
		v.Add("name", r.Name)
	}

	if r.Hydrated {
		v.Add("hydrated", "true")
	}

	if r.Archived != nil {
		v.Add("archived", boolString[*r.Archived])
	}

	if len(r.Clients) > 0 {
		v.Add("clients", strings.Join(r.Clients, ","))
	}

	u.RawQuery = v.Encode()

	return u
}

// GetProjectRequest query parameters to fetch a project
type GetProjectRequest struct {
	Hydrated bool
}

// AppendToQuery decorates the URL with a query string
func (r GetProjectRequest) AppendToQuery(u *url.URL) *url.URL {

	v := u.Query()

	if r.Hydrated {
		v.Add("hydrated", "true")
	}

	u.RawQuery = v.Encode()
	return u
}

// AddProjectRequest represents the parameters to create a project
type AddProjectRequest struct {
	Name     string `json:"name"`
	ClientId string `json:"clientId,omitempty"`
	IsPublic bool   `json:"isPublic"`
	Color    string `json:"color,omitempty"`
	Note     string `json:"note,omitempty"`
	Billable bool   `json:"billable"`
	Public   bool   `json:"public"`
}

// UpdateProjectRequest represents the parameters to update a project
type UpdateProjectRequest struct {
	Name     *string `json:"name,omitempty"`
	ClientId *string `json:"clientId,omitempty"`
	IsPublic *bool   `json:"isPublic,omitempty"`
	Color    *string `json:"color,omitempty"`
	Note     *string `json:"note,omitempty"`
	Billable *bool   `json:"billable,omitempty"`
	Archived *bool   `json:"archived,omitempty"`
}

// UpdateProjectMembershipsRequest represents a request to change which users
// and groups have access to a project
type UpdateProjectMembershipsRequest struct {
	Memberships []UpdateProjectMembership `json:"memberships"`
}

// UpdateProjectMembership sets which user or group has access, and their
// hourly rate
type UpdateProjectMembership struct {
	UserID     string `json:"userId"`
	HourlyRate Rate   `json:"hourlyRate"`
}

// UpdateProjectTemplateRequest represents a request to change isTemplate flag
// of a project
type UpdateProjectTemplateRequest struct {
	IsTemplate bool `json:"isTemplate"`
}

type GetTagsRequest struct {
	Name     string
	Archived *bool

	pagination
}

// WithPagination add pagination to the GetTagsRequest
func (r GetTagsRequest) WithPagination(page, size int) PaginatedRequest {
	r.pagination = newPagination(page, size)
	return r
}

// AppendToQuery decorates the URL with the query string needed for this Request
func (r GetTagsRequest) AppendToQuery(u *url.URL) *url.URL {
	u = r.pagination.AppendToQuery(u)

	v := u.Query()
	if r.Name != "" {
		v.Add("name", r.Name)
	}

	if r.Archived != nil {
		v.Add("archived", boolString[*r.Archived])
	}

	u.RawQuery = v.Encode()

	return u
}

// GetTasksRequest represents the query filters to search tasks of a project
type GetTasksRequest struct {
	Name   string
	Active bool

	pagination
}

// WithPagination add pagination to the GetTasksRequest
func (r GetTasksRequest) WithPagination(page, size int) PaginatedRequest {
	r.pagination = newPagination(page, size)
	return r
}

// AppendToQuery decorates the URL with the query string needed for this Request
func (r GetTasksRequest) AppendToQuery(u *url.URL) *url.URL {
	u = r.pagination.AppendToQuery(u)

	v := u.Query()
	if r.Name != "" {
		v.Add("name", r.Name)
	}

	if r.Active {
		v.Add("is-active", "true")
	}

	u.RawQuery = v.Encode()

	return u
}

type AddTaskRequest struct {
	Name        string    `json:"name"`
	AssigneeIDs *[]string `json:"assigneeIds,omitempty"`
	Billable    *bool     `json:"billable,omitempty"`
	Estimate    *Duration `json:"estimate,omitempty"`
	Status      *string   `json:"status,omitempty"`
}

type UpdateTaskRequest struct {
	Name        string    `json:"name"`
	AssigneeIDs *[]string `json:"assigneeIds,omitempty"`
	Billable    *bool     `json:"billable,omitempty"`
	Estimate    *Duration `json:"estimate,omitempty"`
	Status      *string   `json:"status,omitempty"`
}

type ChangeTimeEntriesInvoicedRequest struct {
	TimeEntryIDs []string `json:"timeEntryIds"`
	Invoiced     bool     `json:"invoiced"`
}

type WorkspaceUsersRequest struct {
	Email string
	pagination
}

// WithPagination add pagination to the WorkspaceUsersRequest
func (r WorkspaceUsersRequest) WithPagination(page, size int) PaginatedRequest {
	r.pagination = newPagination(page, size)
	return r
}

// AppendToQuery decorates the URL with the query string needed for this Request
func (r WorkspaceUsersRequest) AppendToQuery(u *url.URL) *url.URL {
	u = r.pagination.AppendToQuery(u)

	v := u.Query()

	if r.Email != "" {
		v.Add("email", r.Email)
	}

	u.RawQuery = v.Encode()

	return u
}

// UpdateProjectUserRateRequest represents a request to change a user
// billable rate on a project
type UpdateProjectUserRateRequest struct {
	Amount uint      `json:"amount"`
	Since  *DateTime `json:"since,omitempty"`
}

// BaseEstimateRequest is basic information to estime a project
type BaseEstimateRequest struct {
	Type         *EstimateType        `json:"type,omitempty"`
	Active       bool                 `json:"active"`
	ResetOptions *EstimateResetOption `json:"resetOption,omitempty"`
}

// TimeEstimateRequest set parameters for time estimate on a project
type TimeEstimateRequest struct {
	BaseEstimateRequest
	Estimate *Duration `json:"estimate,omitempty"`
}

// BudgetEstimateRequest set parameters for time estimate on a project
type BudgetEstimateRequest struct {
	BaseEstimateRequest
	Estimate *uint64 `json:"estimate,omitempty"`
}

// UpdateProjectEstimateRequest represents a request to set a estimate of a
// project
type UpdateProjectEstimateRequest struct {
	TimeEstimate   TimeEstimateRequest   `json:"timeEstimate"`
	BudgetEstimate BudgetEstimateRequest `json:"budgetEstimate"`
}
