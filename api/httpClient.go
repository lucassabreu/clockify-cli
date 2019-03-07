package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// ErrorNotFound Not Found
var ErrorNotFound = dto.Error{Message: "Nothing was found"}

type transport struct {
	apiKey string
	next   http.RoundTripper
}

func (t transport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("X-Api-Key", t.apiKey)

	return t.next.RoundTrip(r)
}

// NewRequest to be used in Client
func (c *Client) NewRequest(method, uri string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(strings.Join([]string{c.baseURL.Path, uri}, "/"))
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
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
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	r, err := c.Client.Do(req)
	if err != nil {
		return r, err
	}

	if v == nil {
		return r, err
	}

	if r.StatusCode == 404 {
		return r, ErrorNotFound
	}

	defer r.Body.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, r.Body)
	c.debugf("url: %s, status: %d, body: %s", req.URL.String(), r.StatusCode, buf)

	decoder := json.NewDecoder(buf)

	if r.StatusCode < 200 || r.StatusCode > 300 {
		var apiErr dto.Error
		err = decoder.Decode(&apiErr)
		if err != nil {
			return r, err
		}
		return r, apiErr
	}

	err = decoder.Decode(v)
	return r, err
}
