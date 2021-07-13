package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	stackedErrors "github.com/pkg/errors"
)

// ErrorMissingAPIKey returned if X-Api-Key is missing
var ErrorMissingAPIKey = errors.New("api Key must be informed")

type Client struct {
	baseURL url.URL
	http.Client
	Logger Logger
}

type transport struct {
	apiKey string
	next   http.RoundTripper
}

func (t transport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("X-Api-Key", t.apiKey)

	return t.next.RoundTrip(r)
}

// NewHttpClient create a new Client, based on: https://clockify.me/developers-api
func NewHttpClient(baseURL, apiKey string) (*Client, error) {
	if len(apiKey) == 0 {
		return nil, stackedErrors.WithStack(ErrorMissingAPIKey)
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, stackedErrors.WithStack(err)
	}

	c := &Client{
		baseURL: *u,
		Client: http.Client{
			Transport: transport{
				apiKey: apiKey,
				next:   http.DefaultTransport,
			},
		},
	}

	return c, nil
}

type Logger interface {
	Printf(string, ...interface{})
}

func (c *Client) logf(format string, v ...interface{}) {
	if c.Logger == nil {
		return
	}

	c.Logger.Printf(format, v)
}

// QueryAppender an interface to identify if the parameters should be sent through the query or body
type QueryAppender interface {
	AppendToQuery(url.URL) url.URL
}

// Error api errors
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (code: %d)", e.Message, e.Code)
}

// ErrorNotFound Not Found
var ErrorNotFound = Error{Message: "Nothing was found"}

// NewRequest to be used in Client
func (c *Client) NewRequest(method, uri string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(strings.Join([]string{c.baseURL.Path, uri}, "/"))
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
		c.logf("request body: %s", buf.(*bytes.Buffer))
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
	defer r.Body.Close()

	buf := new(bytes.Buffer)

	_, err = io.Copy(buf, r.Body)
	if err != nil {
		return nil, stackedErrors.WithStack(err)
	}

	c.logf("url: %s, status: %d, body: \"%s\"", req.URL.String(), r.StatusCode, buf)

	if r.StatusCode == 404 {
		return r, stackedErrors.WithStack(ErrorNotFound)
	}

	decoder := json.NewDecoder(buf)

	if r.StatusCode < 200 || r.StatusCode > 300 {
		var apiErr Error
		err = decoder.Decode(&apiErr)
		if err != nil {
			return r, stackedErrors.WithStack(err)
		}
		return r, stackedErrors.WithStack(apiErr)
	}

	if v == nil {
		return r, nil
	}

	if buf.Len() == 0 {
		return r, nil
	}

	return r, stackedErrors.WithStack(decoder.Decode(v))
}
