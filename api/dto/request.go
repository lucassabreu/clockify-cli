package dto

import (
	"encoding/json"
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

	if len(s) < 4 {
		return errors.Errorf("duration %s is invalid", b)
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
			return errors.Wrap(err, "unmarshal duration")
		}
		dc = dc + time.Duration(v)*u
		j = i + 1
	}

	*d = Duration{Duration: dc}
	return nil
}

func (d Duration) String() string {
	return "PT" + strings.ToUpper(d.Duration.String())
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

// TimeEntryStartEndRequest to get a time entry
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
	Start        DateTime      `json:"start,omitempty"`
	End          *DateTime     `json:"end,omitempty"`
	Billable     bool          `json:"billable,omitempty"`
	Description  string        `json:"description,omitempty"`
	ProjectID    string        `json:"projectId,omitempty"`
	TaskID       string        `json:"taskId,omitempty"`
	TagIDs       []string      `json:"tagIds,omitempty"`
	CustomFields []CustomField `json:"customFields,omitempty"`
}

// CustomField DTO
type CustomField struct {
	CustomFieldID string `json:"customFieldId"`
	Value         string `json:"value"`
}

// UpdateTimeEntryRequest to update a time entry
type UpdateTimeEntryRequest struct {
	Start        DateTime      `json:"start,omitempty"`
	End          *DateTime     `json:"end,omitempty"`
	Billable     bool          `json:"billable,omitempty"`
	Description  string        `json:"description,omitempty"`
	ProjectID    string        `json:"projectId,omitempty"`
	TaskID       string        `json:"taskId,omitempty"`
	TagIDs       []string      `json:"tagIds,omitempty"`
	CustomFields []CustomField `json:"customFields,omitempty"`
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

type GetProjectRequest struct {
	Name     string
	Archived *bool
	Clients  []string

	pagination
}

// WithPagination add pagination to the GetProjectRequest
func (r GetProjectRequest) WithPagination(page, size int) PaginatedRequest {
	r.pagination = newPagination(page, size)
	return r
}

var boolString = map[bool]string{
	true:  "true",
	false: "false",
}

// AppendToQuery decorates the URL with the query string needed for this Request
func (r GetProjectRequest) AppendToQuery(u *url.URL) *url.URL {
	u = r.pagination.AppendToQuery(u)

	v := u.Query()

	if r.Name != "" {
		v.Add("name", r.Name)
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
	v.Add("name", r.Name)
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
	v.Add("name", r.Name)
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
