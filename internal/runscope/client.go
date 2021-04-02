package runscope

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const DefaultEndpoint = "https://api.runscope.com"

type Client struct {
	endpoint   string
	token      string
	httpClient *http.Client

	Test        TestClient
	Environment EnvironmentClient
	Bucket      BucketClient
	Integration IntegrationClient
	Schedule    ScheduleClient
	Step        StepClient
}

func NewClient(options ...ClientOption) *Client {
	client := &Client{
		endpoint:   DefaultEndpoint,
		httpClient: &http.Client{},
	}

	for _, option := range options {
		option(client)
	}

	client.Test = TestClient{client: client}
	client.Environment = EnvironmentClient{client: client}
	client.Bucket = BucketClient{client: client}
	client.Integration = IntegrationClient{client: client}
	client.Schedule = ScheduleClient{client: client}
	client.Step = StepClient{client: client}

	return client
}

type ClientOption func(*Client)

func WithToken(token string) ClientOption {
	return func(client *Client) {
		client.token = token
	}
}

func WithEndpoint(endpoint string) ClientOption {
	return func(client *Client) {
		client.endpoint = strings.TrimRight(endpoint, "/")
	}
}

func (c *Client) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	apiUrl := c.endpoint + path

	req, err := func() (*http.Request, error) {
		if body == nil {
			return http.NewRequest(method, apiUrl, nil)
		}

		reqBodyData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		return http.NewRequest(method, apiUrl, bytes.NewReader(reqBodyData))
	}()

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req = req.WithContext(ctx)

	return req, nil
}

func (c *Client) Do(r *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(r)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		err = &Error{
			Response: resp,
		}
		json.Unmarshal(body, err)
		return err
	}

	if v != nil {
		if err := json.Unmarshal(body, v); err != nil {
			return err
		}
	}

	return nil
}
