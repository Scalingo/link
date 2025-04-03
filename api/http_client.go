package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Scalingo/go-utils/errors/v2"
)

var _ Client = HTTPClient{}

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
	req, err := c.getRequest(ctx, http.MethodGet, "/version", nil)
	if err != nil {
		return "", err
	}

	resp, err := c.getClient().Do(req)
	if err != nil {
		return "", errors.Wrap(ctx, err, "do request")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", getErrorFromBody(ctx, resp.StatusCode, resp.Body)
	}

	var res map[string]string
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", errors.Wrap(ctx, err, "read version JSON")
	}
	return res["version"], nil
}

func (c HTTPClient) ListEndpoints(ctx context.Context) ([]Endpoint, error) {
	req, err := c.getRequest(ctx, http.MethodGet, "/endpoints", nil)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create request")
	}

	resp, err := c.getClient().Do(req)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "do request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, getErrorFromBody(ctx, resp.StatusCode, resp.Body)
	}

	res := EndpointListResponse{}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "read endpoints JSON")
	}

	return res.Endpoints, nil
}

func (c HTTPClient) GetEndpoint(ctx context.Context, id string) (Endpoint, error) {
	req, err := c.getRequest(ctx, http.MethodGet, "/endpoints/"+id, nil)
	if err != nil {
		return Endpoint{}, errors.Wrap(ctx, err, "create request")
	}

	resp, err := c.getClient().Do(req)
	if err != nil {
		return Endpoint{}, errors.Wrap(ctx, err, "do request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Endpoint{}, getErrorFromBody(ctx, resp.StatusCode, resp.Body)
	}

	res := EndpointGetResponse{}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return Endpoint{}, errors.Wrap(ctx, err, "read endpoint JSON")
	}

	return res.Endpoint, nil
}

func (c HTTPClient) Failover(ctx context.Context, id string) error {
	req, err := c.getRequest(ctx, http.MethodPost, fmt.Sprintf("/endpoints/%s/failover", id), nil)
	if err != nil {
		return err
	}

	resp, err := c.getClient().Do(req)
	if err != nil {
		return errors.Wrap(ctx, err, "do request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return getErrorFromBody(ctx, resp.StatusCode, resp.Body)
	}

	return nil
}

func (c HTTPClient) AddEndpoint(ctx context.Context, params AddEndpointParams) (Endpoint, error) {
	buffer := &bytes.Buffer{}
	err := json.NewEncoder(buffer).Encode(params)
	if err != nil {
		return Endpoint{}, errors.Wrap(ctx, err, "encode endpoint")
	}

	req, err := c.getRequest(ctx, http.MethodPost, "/endpoints", buffer)
	if err != nil {
		return Endpoint{}, errors.Wrap(ctx, err, "create request")
	}

	resp, err := c.getClient().Do(req)
	if err != nil {
		return Endpoint{}, errors.Wrap(ctx, err, "do request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Endpoint{}, getErrorFromBody(ctx, resp.StatusCode, resp.Body)
	}

	res := Endpoint{}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return Endpoint{}, errors.Wrap(ctx, err, "read endpoint JSON")
	}

	return res, nil
}

func (c HTTPClient) UpdateEndpoint(ctx context.Context, id string, params UpdateEndpointParams) (Endpoint, error) {
	buffer := &bytes.Buffer{}
	err := json.NewEncoder(buffer).Encode(params)
	if err != nil {
		return Endpoint{}, errors.Wrap(ctx, err, "encode endpoint")
	}

	req, err := c.getRequest(ctx, http.MethodPatch, "/endpoints/"+id, buffer)
	if err != nil {
		return Endpoint{}, errors.Wrap(ctx, err, "create request")
	}

	resp, err := c.getClient().Do(req)
	if err != nil {
		return Endpoint{}, errors.Wrap(ctx, err, "do request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Endpoint{}, getErrorFromBody(ctx, resp.StatusCode, resp.Body)
	}

	var ip Endpoint
	err = json.NewDecoder(resp.Body).Decode(&ip)
	if err != nil {
		return Endpoint{}, errors.Wrap(ctx, err, "read endpoint JSON")
	}

	return ip, nil
}

func (c HTTPClient) RemoveEndpoint(ctx context.Context, id string) error {
	req, err := c.getRequest(ctx, http.MethodDelete, "/endpoints/"+id, nil)
	if err != nil {
		return errors.Wrap(ctx, err, "create request")
	}

	resp, err := c.getClient().Do(req)
	if err != nil {
		return errors.Wrap(ctx, err, "do request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return getErrorFromBody(ctx, resp.StatusCode, resp.Body)
	}

	return nil
}

func (c HTTPClient) getClient() *http.Client {
	return &http.Client{
		Timeout: c.timeout,
	}
}

func (c HTTPClient) getRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	var url string
	if path[0] != '/' {
		path = "/" + path
	}
	url = fmt.Sprintf("http://%s%s", c.url, path)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create request")
	}

	if c.username != "" || c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	return req, nil
}

func getErrorFromBody(ctx context.Context, statusCode int, body io.Reader) error {
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
			return errors.New(ctx, value)
		}
	}

	return fmt.Errorf("Unexpected status code: %v", statusCode)
}
