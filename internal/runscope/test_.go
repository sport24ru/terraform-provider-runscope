package runscope

import (
	"context"
	"fmt"
	"time"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope/schema"
)

type TestMinimal struct {
	Name        string
	Description string
}

func (opts *TestMinimal) setRequest(body *schema.TestMinimal) {
	body.Name = opts.Name
	body.Description = opts.Description
}

type Test struct {
	TestMinimal
	Id                   string
	DefaultEnvironmentId string
	Steps                []TestStep
	CreatedAt            time.Time
	CreatedBy            CreatedBy
	LastRun              time.Time
	TriggerURL           string
}

type TestStep struct {
	Id string
}

type CreatedBy struct {
	Id    string
	Name  string
	Email string
}

type TestClient struct {
	client *Client
}

func TestFromSchema(s schema.Test) *Test {
	test := &Test{}
	test.Id = s.Id
	test.Name = s.Name
	test.Description = s.Description
	test.DefaultEnvironmentId = s.DefaultEnvironmentId
	test.Steps = make([]TestStep, len(s.Steps))
	for i, step := range s.Steps {
		test.Steps[i].Id = step.Id
	}
	test.CreatedAt = time.Unix(s.CreatedAt, 0)
	test.CreatedBy = CreatedBy{
		Id:    s.CreatedBy.Id,
		Name:  s.CreatedBy.Name,
		Email: s.CreatedBy.Email,
	}
	test.TriggerURL = s.TriggerURL
	return test
}

type TestGetOpts struct {
	BucketId string
	Id       string
}

func (c *TestClient) Get(ctx context.Context, opts TestGetOpts) (*Test, error) {
	req, err := c.client.NewRequest(ctx,
		"GET", fmt.Sprintf("/buckets/%s/tests/%s", opts.BucketId, opts.Id),
		nil)
	if err != nil {
		return nil, err
	}

	var resp schema.TestGetResponse
	err = c.client.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	return TestFromSchema(resp.Test), err
}

type TestCreateOpts struct {
	TestMinimal
	BucketId string
}

func (c *TestClient) Create(ctx context.Context, opts TestCreateOpts) (*Test, error) {
	body := schema.TestCreateRequest{}
	opts.TestMinimal.setRequest(&body.TestMinimal)

	req, err := c.client.NewRequest(ctx,
		"POST", fmt.Sprintf("/buckets/%s/tests", opts.BucketId),
		&body)
	if err != nil {
		return nil, err
	}

	var resp schema.TestCreateResponse
	err = c.client.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	return TestFromSchema(resp.Test), err
}

type TestUpdateOpts struct {
	Test
	BucketId string
}

func (opts *TestUpdateOpts) setRequest(body *schema.TestUpdateRequest) {
	body.Name = opts.Name
	body.Description = opts.Description
	body.DefaultEnvironmentId = opts.DefaultEnvironmentId
}

func (c *TestClient) Update(ctx context.Context, opts TestUpdateOpts) (*Test, error) {
	body := schema.TestUpdateRequest{}
	opts.setRequest(&body)

	req, err := c.client.NewRequest(ctx,
		"PUT", fmt.Sprintf("/buckets/%s/tests/%s", opts.BucketId, opts.Id),
		&body)
	if err != nil {
		return nil, err
	}

	var resp schema.TestUpdateResponse
	err = c.client.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	return TestFromSchema(resp.Test), err
}

type TestDeleteOpts struct {
	BucketId string
	Id       string
}

func (c *TestClient) Delete(ctx context.Context, opts TestDeleteOpts) error {
	req, err := c.client.NewRequest(ctx,
		"DELETE", fmt.Sprintf("/buckets/%s/tests/%s", opts.BucketId, opts.Id),
		nil)
	if err != nil {
		return err
	}

	err = c.client.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}
