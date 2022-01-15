package dto

import (
	"net/url"
	"strconv"
	"time"
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
func (p pagination) AppendToQuery(u url.URL) url.URL {
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
func (r GetTimeEntryRequest) AppendToQuery(u url.URL) url.URL {
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
func (r UserTimeEntriesRequest) AppendToQuery(u url.URL) url.URL {
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

type GetProjectRequest struct {
	Name     string
	Archived *bool

	pagination
}

// WithPagination add pagination to the GetProjectRequest
func (r GetProjectRequest) WithPagination(page, size int) PaginatedRequest {
	r.pagination = newPagination(page, size)
	return r
}

func bool2str(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// AppendToQuery decorates the URL with the query string needed for this Request
func (r GetProjectRequest) AppendToQuery(u url.URL) url.URL {
	u = r.pagination.AppendToQuery(u)

	v := u.Query()

	if r.Name != "" {
		v.Add("name", r.Name)
	}

	if r.Archived != nil {
		v.Add("archived", bool2str(*r.Archived))
	}

	u.RawQuery = v.Encode()

	return u
}

type GetTagsRequest struct {
	Name     string
	Archived bool

	pagination
}

// WithPagination add pagination to the GetTagsRequest
func (r GetTagsRequest) WithPagination(page, size int) PaginatedRequest {
	r.pagination = newPagination(page, size)
	return r
}

// AppendToQuery decorates the URL with the query string needed for this Request
func (r GetTagsRequest) AppendToQuery(u url.URL) url.URL {
	u = r.pagination.AppendToQuery(u)

	v := u.Query()
	v.Add("name", r.Name)
	if r.Archived {
		v.Add("archived", "true")
	}

	u.RawQuery = v.Encode()

	return u
}

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
func (r GetTasksRequest) AppendToQuery(u url.URL) url.URL {
	u = r.pagination.AppendToQuery(u)

	v := u.Query()
	v.Add("name", r.Name)
	if r.Active {
		v.Add("active", "true")
	}

	u.RawQuery = v.Encode()

	return u
}

type ChangeTimeEntriesInvoicedRequest struct {
	TimeEntryIDs []string `json:"timeEntryIds"`
	Invoiced     bool     `json:"invoiced"`
}
