package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/pkg/errors"
)

// QueryAppender an interface to identify if the parameters should be sent through the query or body
type QueryAppender interface {
	AppendToQuery(url.URL) url.URL
}

// ErrorNotFound Not Found
var ErrorNotFound = dto.Error{Message: "Nothing was found", Code: 404}

// ErrorForbidden Forbidden
var ErrorForbidden = dto.Error{Message: "Forbidden", Code: 403}

type transport struct {
	apiKey string
	next   http.RoundTripper
}

func (t transport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("X-Api-Key", t.apiKey)

	return t.next.RoundTrip(r)
}

// NewRequest to be used in Client
func (c *client) NewRequest(method, uri string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(c.baseURL.Path + "/" + uri)
	if err != nil {
		return nil, err
	}

	if qa, ok := body.(QueryAppender); ok {
		*u = qa.AppendToQuery(*u)
	}

	if method == "GET" {
		body = nil
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
		c.infof("request body: %s", buf.(*bytes.Buffer))
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	return req, nil
}

// Do executes a http.Request inside the Clockify's Client
func (c *client) Do(
	req *http.Request, v interface{}, name string) (*http.Response, error) {
	r, err := c.Client.Do(req)
	if err != nil {
		return r, err
	}
	defer r.Body.Close()

	buf := new(bytes.Buffer)

	_, err = io.Copy(buf, r.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if c.debugLogger != nil {
		c.debugf("name: %s, method: %s, url: %s, status: %d, response: \"%s\"",
			name, req.Method, req.URL.String(), r.StatusCode, buf)
	} else {
		c.infof("name: %s, method: %s, url: %s, status: %d",
			name, req.Method, req.URL.String(), r.StatusCode)
	}

	decoder := json.NewDecoder(buf)

	if r.StatusCode < 200 || r.StatusCode > 300 {
		var apiErr dto.Error
		err = decoder.Decode(&apiErr)
		if err != nil && err != io.EOF {
			return r, errors.WithStack(err)
		}

		if r.StatusCode == 404 && apiErr.Message == "" {
			apiErr = ErrorNotFound
		}

		if r.StatusCode == 403 && apiErr.Message == "" {
			apiErr = ErrorForbidden
		}

		if apiErr.Message == "" {
			apiErr.Message = "No response"
		}

		return r, errors.WithStack(apiErr)
	}

	if v == nil {
		return r, nil
	}

	if buf.Len() == 0 {
		return r, nil
	}

	return r, errors.WithStack(decoder.Decode(v))
}
