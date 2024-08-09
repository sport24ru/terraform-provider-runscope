package runscope

import (
	"context"
	"fmt"
	"net/url"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope/schema"
)

type Bucket struct {
	Key        string
	Name       string
	Team       Team
	AuthToken  string
	Default    bool
	VerifySSL  bool
	TriggerURL string
}

type Team struct {
	Name string
	UUID string
}

type BucketClient struct {
	client *Client
}

func BucketFromSchema(s *schema.Bucket) *Bucket {
	return &Bucket{
		Key:  s.Key,
		Name: s.Name,
		Team: Team{
			Name: s.Team.Name,
			UUID: s.Team.Id,
		},
		AuthToken:  s.AuthToken,
		Default:    s.Default,
		VerifySSL:  s.VerifySSL,
		TriggerURL: s.TriggerURL,
	}
}

const bucketsBaseUrl = "/buckets"

type BucketCreateOpts struct {
	Name     string
	TeamUUID string
}

func (opts *BucketCreateOpts) URL() string {
	return fmt.Sprintf("%s?name=%s&team_uuid=%s",
		bucketsBaseUrl, url.QueryEscape(opts.Name), url.QueryEscape(opts.TeamUUID))
}

func (c *BucketClient) Create(ctx context.Context, opts *BucketCreateOpts) (*Bucket, error) {
	req, err := c.client.NewRequest(ctx, "POST", opts.URL(), nil)
	if err != nil {
		return nil, err
	}

	var resp schema.BucketCreateResponse
	err = c.client.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	return BucketFromSchema(&resp.Bucket), err
}

type BucketGetOpts struct {
	Key string
}

func (opts *BucketGetOpts) URL() string {
	return fmt.Sprintf("/buckets/%s", opts.Key)
}

func (c *BucketClient) Get(ctx context.Context, opts *BucketGetOpts) (*Bucket, error) {
	req, err := c.client.NewRequest(ctx, "GET", opts.URL(), nil)
	if err != nil {
		return nil, err
	}

	var resp schema.BucketGetResponse
	err = c.client.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	return BucketFromSchema(&resp.Bucket), nil
}

func (c *BucketClient) List(ctx context.Context) ([]*Bucket, error) {
	req, err := c.client.NewRequest(ctx, "GET", bucketsBaseUrl, nil)
	if err != nil {
		return nil, err
	}

	var resp schema.BucketListResponse
	err = c.client.Do(req, &resp)
	if err != nil {
		return nil, err
	}
	buckets := make([]*Bucket, len(resp.Buckets), len(resp.Buckets))
	for i, bucket := range resp.Buckets {
		buckets[i] = BucketFromSchema(&bucket)
	}

	return buckets, nil
}

type BucketDeleteOpts struct {
	BucketGetOpts
}

func (c *BucketClient) Delete(ctx context.Context, opts *BucketDeleteOpts) error {
	req, err := c.client.NewRequest(ctx, "DELETE", opts.URL(), nil)
	if err != nil {
		return err
	}

	return c.client.Do(req, nil)
}
