package httpreq

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Client struct {
	url     string
	headers map[string]string
	body    any
	method  string
	client  *http.Client
}

func NewClient(timeout ...time.Duration) *Client {
	var t time.Duration
	if len(timeout) > 0 {
		t = timeout[0]
	} else {
		t = 10 * time.Second
	}
	return &Client{
		client: &http.Client{
			Timeout: t,
		},
	}
}

func (c *Client) Url(url string) *Client {
	c.url = url
	return c
}

func (c *Client) Headers(headers map[string]string) *Client {
	c.headers = headers
	return c
}

func (c *Client) Body(body any) *Client {
	c.body = body
	return c
}

func (c *Client) Get() (*http.Response, error) {
	c.method = http.MethodGet
	return c.send()
}

func (c *Client) Post() (*http.Response, error) {
	c.method = http.MethodPost
	return c.send()
}

func (c *Client) Put() (*http.Response, error) {
	c.method = http.MethodPut
	return c.send()
}

func (c *Client) Delete() (*http.Response, error) {
	c.method = http.MethodDelete
	return c.send()
}

func (c *Client) send() (*http.Response, error) {
	var body io.Reader
	if c.body != nil {
		jsonBody, err := json.Marshal(c.body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonBody)
	}
	req, err := http.NewRequest(c.method, c.url, body)
	if err != nil {
		return nil, err
	}
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	return c.client.Do(req)
}
