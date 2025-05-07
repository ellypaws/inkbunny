package inkbunny

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ellypaws/inkbunny/utils"
)

var DefaultClient = NewClient()

type Client struct {
	ctx    context.Context
	client *http.Client
}

func (c *Client) Get() *Client {
	if c != nil {
		return c
	}
	return DefaultClient
}

func NewClient(opts ...func(*Client)) *Client {
	c := new(Client)
	for _, opt := range opts {
		opt(c)
	}
	if c.ctx == nil {
		c.ctx = context.Background()
	}
	if c.client == nil {
		c.client = &http.Client{
			Timeout: 5 * time.Minute,
		}
	}
	return c
}

func WithClient(client *http.Client) func(*Client) {
	return func(c *Client) {
		c.client = client
	}
}

func WithContext(ctx context.Context) func(*Client) {
	return func(c *Client) {
		c.ctx = ctx
	}
}

func (c *Client) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func (c *Client) SetClient(client *http.Client) {
	c.client = client
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.Timeout = timeout
}

const (
	MimeTypeJSON  = "application/json"
	MimeTypeForm  = "multipart/form-data"
	MimeTypeQuery = "application/x-www-form-urlencoded"
)

// PostDecode sends a POST request to the given URL with the provided data.
// It automatically reads the [http.Response.Body], checks for errors and decodes into T.
// It calls Client.PostForm and then utils.ParseResponse[T].
func PostDecode[T any](c *Client, url *url.URL, data any) (T, error) {
	response, err := c.PostForm(url, data)
	if err != nil {
		var t T
		return t, err
	}
	return utils.ParseResponse[T](response)
}

// PostForm sends a POST request to the specified URL with the provided data and returns the HTTP response or an error.
// The method determines the appropriate content type and body format based on the type of the data parameter.
// Passing in a []byte or any type that implements io.Reader assumes the Content-Type is of MimeTypeJSON.
func (c *Client) PostForm(u *url.URL, data any) (*http.Response, error) {
	contentType := MimeTypeQuery
	var body io.Reader
	switch d := data.(type) {
	case nil:
	case []byte:
		body = bytes.NewReader(d)
		contentType = MimeTypeJSON
	case io.Reader:
		body = d
		contentType = MimeTypeJSON
	case url.Values:
		if u.RawQuery != "" {
			u.RawQuery += "&" + d.Encode()
		} else {
			u.RawQuery = d.Encode()
		}
		body = strings.NewReader(u.RawQuery)
	default:
		var err error
		body, contentType, err = utils.StructToMultipartForm(data)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(c.ctx, http.MethodPost, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	return c.client.Do(req)
}
