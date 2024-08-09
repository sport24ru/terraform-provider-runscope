package runscope

import (
	"context"
	"fmt"
	"time"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope/schema"
)

type ScheduleBase struct {
	EnvironmentId string
	Interval      string
	Note          string
}

type Schedule struct {
	ScheduleBase
	Id         string
	ExportedAt time.Time
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
	schedule.ExportedAt = time.Unix(s.ExportedAt, 0)
	return schedule
}

type ScheduleURLOpts struct {
	BucketId string
	TestId   string
}

func (opts *ScheduleURLOpts) URL() string {
	return fmt.Sprintf("/buckets/%s/tests/%s/schedules", opts.BucketId, opts.TestId)
}

type ScheduleCreateOpts struct {
	ScheduleURLOpts
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
	ScheduleURLOpts
	Id string
}

func (opts *ScheduleGetOpts) URL() string {
	return fmt.Sprintf("%s/%s", opts.ScheduleURLOpts.URL(), opts.Id)
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

type ScheduleUpdateOpts struct {
	ScheduleGetOpts
	ScheduleBase
}

func (opts *ScheduleUpdateOpts) setRequest(body *schema.ScheduleUpdateRequest) {
	body.Note = opts.Note
	body.Interval = opts.Interval
	body.EnvironmentId = opts.EnvironmentId
}

func (c *ScheduleClient) Update(ctx context.Context, opts *ScheduleUpdateOpts) (*Schedule, error) {
	body := schema.ScheduleUpdateRequest{}
	opts.setRequest(&body)

	req, err := c.client.NewRequest(ctx, "PUT", opts.URL(), &body)
	if err != nil {
		return nil, err
	}

	var resp schema.ScheduleUpdateResponse
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
