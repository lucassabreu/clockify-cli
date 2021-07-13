package dto

import (
	"net/url"
	"strconv"

	"github.com/lucassabreu/clockify-cli/http"
)

type pagination struct {
	page     int
	pageSize int
}

func NewPagination(page, size int) pagination {
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

// TimeEntryStartEndRequest to get entries by range
type TimeEntryStartEndRequest struct {
	Start    http.DateTime
	End      http.DateTime
	Hydrated *bool

	Pagination pagination
}

// WithPagination add pagination to the TimeEntryStartEndRequest
func (r TimeEntryStartEndRequest) WithPagination(page, size int) PaginatedRequest {
	r.Pagination = NewPagination(page, size)
	return r
}

// AppendToQuery decorates the URL with the query string needed for this Request
func (r TimeEntryStartEndRequest) AppendToQuery(u url.URL) url.URL {
	u = r.Pagination.AppendToQuery(u)
	v := u.Query()
	v.Add("start", r.Start.String())
	v.Add("end", r.End.String())
	if r.Hydrated != nil && *r.Hydrated {
		v.Add("hydrated", "true")
	}

	u.RawQuery = v.Encode()

	return u
}

// OutTimeEntryRequest to end the current time entry
type OutTimeEntryRequest struct {
	End http.DateTime `json:"end"`
}

// CreateTimeEntryRequest to create a time entry is created
type CreateTimeEntryRequest struct {
	Start       http.DateTime  `json:"start,omitempty"`
	End         *http.DateTime `json:"end,omitempty"`
	Billable    bool           `json:"billable,omitempty"`
	Description string         `json:"description,omitempty"`
	ProjectID   string         `json:"projectId,omitempty"`
	TaskID      string         `json:"taskId,omitempty"`
	TagIDs      []string       `json:"tagIds,omitempty"`
}

// UpdateTimeEntryRequest to update a time entry
type UpdateTimeEntryRequest struct {
	Start       http.DateTime  `json:"start,omitempty"`
	End         *http.DateTime `json:"end,omitempty"`
	Billable    bool           `json:"billable,omitempty"`
	Description string         `json:"description,omitempty"`
	ProjectID   string         `json:"projectId,omitempty"`
	TaskID      string         `json:"taskId,omitempty"`
	TagIDs      []string       `json:"tagIds,omitempty"`
}

type GetProjectRequest struct {
	Name     string
	Archived bool

	Pagination pagination
}

// WithPagination add pagination to the GetProjectRequest
func (r GetProjectRequest) WithPagination(page, size int) PaginatedRequest {
	r.Pagination = NewPagination(page, size)
	return r
}

// AppendToQuery decorates the URL with the query string needed for this Request
func (r GetProjectRequest) AppendToQuery(u url.URL) url.URL {
	u = r.Pagination.AppendToQuery(u)

	v := u.Query()
	v.Add("name", r.Name)
	if r.Archived {
		v.Add("archived", "true")
	}

	u.RawQuery = v.Encode()

	return u
}

type GetTagsRequest struct {
	Name     string
	Archived bool

	Pagination pagination
}

// WithPagination add pagination to the GetTagsRequest
func (r GetTagsRequest) WithPagination(page, size int) PaginatedRequest {
	r.Pagination = NewPagination(page, size)
	return r
}

// AppendToQuery decorates the URL with the query string needed for this Request
func (r GetTagsRequest) AppendToQuery(u url.URL) url.URL {
	u = r.Pagination.AppendToQuery(u)

	v := u.Query()
	v.Add("name", r.Name)
	if r.Archived {
		v.Add("archived", "true")
	}

	u.RawQuery = v.Encode()

	return u
}
