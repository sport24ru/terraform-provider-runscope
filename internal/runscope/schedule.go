package runscope

import (
	"context"
	"fmt"
	"github.com/terraform-providers/terraform-provider-runscope/internal/runscope/schema"
)

type ScheduleBase struct {
	EnvironmentId string
	Interval      string
	Note          string
}

type Schedule struct {
	ScheduleBase
	Id string
}

type ScheduleClient struct {
	client *Client
}

func ScheduleFromSchema(s *schema.Schedule) *Schedule {
	schedule := &Schedule{}
	schedule.Id = s.Id
	schedule.EnvironmentId = s.EnvironmentId
	schedule.Interval = s.Interval
	schedule.Note = s.Note
	return schedule
}

type ScheduleBaseOpts struct {
	BucketId string
	TestId   string
}

func (opts *ScheduleBaseOpts) URL() string {
	return fmt.Sprintf("/buckets/%s/tests/%s/schedules", opts.BucketId, opts.TestId)
}

type ScheduleCreateOpts struct {
	ScheduleBaseOpts
	ScheduleBase
}

func (c *ScheduleClient) Create(ctx context.Context, opts *ScheduleCreateOpts) (*Schedule, error) {
	body := schema.ScheduleCreateRequest{}
	body.EnvironmentId = opts.EnvironmentId
	body.Interval = opts.Interval
	body.Note = opts.Note

	req, err := c.client.NewRequest(ctx, "POST", opts.URL(), &body)
	if err != nil {
		return nil, err
	}

	var resp schema.ScheduleCreateResponse
	err = c.client.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	return ScheduleFromSchema(&resp.Schedule), err
}

type ScheduleGetOpts struct {
	ScheduleBaseOpts
	Id string
}

func (opts *ScheduleGetOpts) URL() string {
	return fmt.Sprintf("%s/%s", opts.ScheduleBaseOpts.URL(), opts.Id)
}

func (c *ScheduleClient) Get(ctx context.Context, opts *ScheduleGetOpts) (*Schedule, error) {
	req, err := c.client.NewRequest(ctx, "GET", opts.URL(), nil)
	if err != nil {
		return nil, err
	}

	var resp schema.ScheduleGetResponse
	err = c.client.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	return ScheduleFromSchema(&resp.Schedule), err
}

type ScheduleDeleteOpts struct {
	ScheduleGetOpts
}

func (c *ScheduleClient) Delete(ctx context.Context, opts *ScheduleDeleteOpts) error {
	req, err := c.client.NewRequest(ctx, "DELETE", opts.URL(), nil)
	if err != nil {
		return err
	}

	err = c.client.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}
