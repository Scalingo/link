package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Scalingo/link/models"
)

type HTTPClient struct {
	url      string
	username string
	password string
	timeout  time.Duration
}

type ClientOpt func(HTTPClient) HTTPClient

func NewHTTPClient(opts ...ClientOpt) HTTPClient {
	client := HTTPClient{}
	for _, opt := range opts {
		client = opt(client)
	}
	return client
}

func WithURL(url string) ClientOpt {
	return func(c HTTPClient) HTTPClient {
		c.url = url
		return c
	}
}
func WithUser(user string) ClientOpt {
	return func(c HTTPClient) HTTPClient {
		c.username = user
		return c
	}
}

func WithPassword(pass string) ClientOpt {
	return func(c HTTPClient) HTTPClient {
		c.password = pass
		return c
	}
}

func WithTimeout(timeout time.Duration) ClientOpt {
	return func(c HTTPClient) HTTPClient {
		c.timeout = timeout
		return c
	}
}

func (c HTTPClient) ListIPs(ctx context.Context) ([]IP, error) {
	req, err := c.getRequest(http.MethodGet, "/ips", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.getClient().Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, getErrorFromBody(resp.StatusCode, resp.Body)
	}

	var res map[string][]IP

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return res["ips"], nil
}

func (c HTTPClient) AddIP(ctx context.Context, ip string, checks ...models.Healthcheck) (IP, error) {
	buffer := &bytes.Buffer{}
	err := json.NewEncoder(buffer).Encode(models.IP{
		IP:     ip,
		Checks: checks,
	})
	if err != nil {
		return IP{}, err
	}

	req, err := c.getRequest(http.MethodPost, "/ips", buffer)
	if err != nil {
		return IP{}, err
	}

	resp, err := c.getClient().Do(req)
	if err != nil {
		return IP{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return IP{}, getErrorFromBody(resp.StatusCode, resp.Body)
	}

	var res IP

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return IP{}, err
	}

	return res, nil
}

func (c HTTPClient) RemoveIP(ctx context.Context, id string) error {
	req, err := c.getRequest(http.MethodDelete, fmt.Sprintf("/ips/%s", id), nil)
	if err != nil {
		return err
	}

	resp, err := c.getClient().Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return getErrorFromBody(resp.StatusCode, resp.Body)
	}

	return nil
}

func (c HTTPClient) getClient() *http.Client {
	return &http.Client{
		Timeout: c.timeout,
	}
}

func (c HTTPClient) getRequest(method, path string, body io.Reader) (*http.Request, error) {
	var url string
	if path[0] != '/' {
		path = "/" + path
	}
	url = fmt.Sprintf("http://%s%s", c.url, path)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if c.username != "" || c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	return req, nil
}

func getErrorFromBody(statusCode int, body io.Reader) error {
	if statusCode/100 == 5 || statusCode/100 == 4 {
		var res map[string]interface{}
		err := json.NewDecoder(body).Decode(&res)
		value := fmt.Sprintf("Unexpected status code: %v", statusCode)
		if err == nil {
			if msg, ok := res["msg"]; ok {
				value = fmt.Sprintf("%v", msg)
			}
			if msg, ok := res["error"]; ok {
				value = fmt.Sprintf("%v", msg)
			}

		}
		return errors.New(value)
	}

	return fmt.Errorf("Unexpected status code: %v", statusCode)
}
