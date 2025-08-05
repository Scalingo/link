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

	"github.com/Scalingo/link/v2/models"
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

func (c HTTPClient) Version(ctx context.Context) (string, error) {
	req, err := c.getRequest(http.MethodGet, "/version", nil)
	if err != nil {
		return "", err
	}

	resp, err := c.getClient().Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", getErrorFromBody(resp.StatusCode, resp.Body)
	}

	var res map[string]string
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", err
	}
	return res["version"], nil
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

func (c HTTPClient) GetIP(ctx context.Context, id string) (IP, error) {
	req, err := c.getRequest(http.MethodGet, fmt.Sprintf("/ips/%s", id), nil)
	if err != nil {
		return IP{}, err
	}

	resp, err := c.getClient().Do(req)
	if err != nil {
		return IP{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return IP{}, getErrorFromBody(resp.StatusCode, resp.Body)
	}

	var res map[string]IP

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return IP{}, err
	}

	return res["ip"], nil
}

func (c HTTPClient) Failover(ctx context.Context, id string) error {
	req, err := c.getRequest(http.MethodPost, fmt.Sprintf("/ips/%s/failover", id), nil)
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

type AddIPParams struct {
	HealthcheckInterval int                  `json:"healthcheck_interval"`
	Checks              []models.Healthcheck `json:"checks"`
}

func (c HTTPClient) AddIP(ctx context.Context, ip string, params AddIPParams) (IP, error) {
	buffer := &bytes.Buffer{}
	err := json.NewEncoder(buffer).Encode(models.IP{
		IP:                  ip,
		HealthcheckInterval: params.HealthcheckInterval,
		Checks:              params.Checks,
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

type UpdateIPParams struct {
	Healthchecks []models.Healthcheck `json:"healthchecks"`
}

func (c HTTPClient) UpdateIP(ctx context.Context, id string, params UpdateIPParams) (IP, error) {
	buffer := &bytes.Buffer{}
	err := json.NewEncoder(buffer).Encode(params)
	if err != nil {
		return IP{}, err
	}

	req, err := c.getRequest(http.MethodPatch, fmt.Sprintf("/ips/%s", id), buffer)
	if err != nil {
		return IP{}, err
	}

	resp, err := c.getClient().Do(req)
	if err != nil {
		return IP{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return IP{}, getErrorFromBody(resp.StatusCode, resp.Body)
	}

	var ip IP
	err = json.NewDecoder(resp.Body).Decode(&ip)
	if err != nil {
		return IP{}, err
	}

	return ip, nil
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
		switch statusCode {
		case 404:
			return ErrNotFound{value}
		default:
			return errors.New(value)
		}
	}

	return fmt.Errorf("Unexpected status code: %v", statusCode)
}
