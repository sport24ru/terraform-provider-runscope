package runscope

import (
	"context"
	"fmt"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope/schema"
)

type EnvironmentBase struct {
	Name                string
	Script              string
	PreserveCookies     bool
	InitialVariables    map[string]string
	Integrations        []string
	Regions             []string
	RemoteAgents        []EnvironmentRemoteAgent
	RetryOnFailure      bool
	StopOnFailure       bool
	VerifySSL           bool
	Webhooks            []string
	Emails              Emails
	ParentEnvironmentId string
	ClientCertificate   string
}

func (eb *EnvironmentBase) setRequest(seb *schema.EnvironmentBase) {
	seb.Name = eb.Name
	seb.Script = eb.Script
	seb.PreserveCookies = eb.PreserveCookies
	seb.InitialVariables = eb.InitialVariables
	seb.Regions = eb.Regions
	seb.RetryOnFailure = eb.RetryOnFailure
	seb.StopOnFailure = eb.StopOnFailure
	seb.VerifySSL = eb.VerifySSL
	seb.Webhooks = eb.Webhooks
	seb.Emails = schema.Emails{
		NotifyAll:       eb.Emails.NotifyAll,
		NotifyOn:        eb.Emails.NotifyOn,
		NotifyThreshold: eb.Emails.NotifyThreshold,
	}
	for _, id := range eb.Integrations {
		seb.Integrations = append(seb.Integrations, schema.EnvironmentIntegration{Id: id})
	}
	for _, agent := range eb.RemoteAgents {
		seb.RemoteAgents = append(seb.RemoteAgents, schema.EnvironmentRemoteAgent{
			Name: agent.Name,
			UUID: agent.UUID,
		})
	}
	seb.Emails.Recipients = make([]schema.Recipient, len(eb.Emails.Recipients))
	for i, recipient := range eb.Emails.Recipients {
		seb.Emails.Recipients[i] = schema.Recipient{
			Id:    recipient.Id,
			Name:  recipient.Name,
			Email: recipient.Email,
		}
	}
	seb.ParentEnvironmentId = eb.ParentEnvironmentId
	seb.ClientCertificate = eb.ClientCertificate
}

type Environment struct {
	EnvironmentBase
	Id string
}

type EnvironmentRemoteAgent struct {
	Name string
	UUID string
}

type Emails struct {
	NotifyAll       bool
	NotifyOn        string
	NotifyThreshold int
	Recipients      []Recipient
}

func (e Emails) IsDefault() bool {
	return !e.NotifyAll && e.NotifyOn == "" && e.NotifyThreshold == 0 && len(e.Recipients) == 0
}

type Recipient struct {
	Id    string
	Name  string
	Email string
}

type EnvironmentClient struct {
	client *Client
}

func EnvironmentFromSchema(s *schema.Environment) *Environment {
	env := &Environment{}
	env.Id = s.Id
	env.Name = s.Name
	env.Script = s.Script
	env.PreserveCookies = s.PreserveCookies
	env.InitialVariables = s.InitialVariables
	env.Regions = s.Regions
	env.RetryOnFailure = s.RetryOnFailure
	env.StopOnFailure = s.StopOnFailure
	env.VerifySSL = s.VerifySSL
	env.Webhooks = s.Webhooks
	env.Emails = Emails{
		NotifyAll:       s.Emails.NotifyAll,
		NotifyOn:        s.Emails.NotifyOn,
		NotifyThreshold: s.Emails.NotifyThreshold,
	}

	for _, i := range s.Integrations {
		env.Integrations = append(env.Integrations, i.Id)
	}
	for _, ra := range s.RemoteAgents {
		env.RemoteAgents = append(env.RemoteAgents, EnvironmentRemoteAgent{
			Name: ra.Name,
			UUID: ra.UUID,
		})
	}
	for _, r := range s.Emails.Recipients {
		env.Emails.Recipients = append(env.Emails.Recipients, Recipient{
			Id:    r.Id,
			Name:  r.Name,
			Email: r.Email,
		})
	}
	env.ParentEnvironmentId = s.ParentEnvironmentId
	env.ClientCertificate = s.ClientCertificate
	return env
}

type EnvironmentUriOpts struct {
	BucketId string
	TestId   string
}

func (opts *EnvironmentUriOpts) BaseURL() string {
	if opts.TestId == "" {
		return fmt.Sprintf("/buckets/%s/environments", opts.BucketId)
	}
	return fmt.Sprintf("/buckets/%s/tests/%s/environments", opts.BucketId, opts.TestId)
}

type EnvironmentCreateOpts struct {
	EnvironmentUriOpts
	EnvironmentBase
}

func (c *EnvironmentClient) Create(ctx context.Context, opts *EnvironmentCreateOpts) (*Environment, error) {
	body := &schema.EnvironmentCreateRequest{}
	opts.EnvironmentBase.setRequest(&body.EnvironmentBase)

	req, err := c.client.NewRequest(ctx, "POST", opts.BaseURL(), &body)
	if err != nil {
		return nil, err
	}

	var resp schema.EnvironmentCreateResponse
	err = c.client.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	return EnvironmentFromSchema(&resp.Environment), err
}

type EnvironmentGetOpts struct {
	EnvironmentUriOpts
	Id string
}

func (opts *EnvironmentGetOpts) URL() string {
	return fmt.Sprintf("%s/%s", opts.BaseURL(), opts.Id)
}

func (c *EnvironmentClient) Get(ctx context.Context, opts *EnvironmentGetOpts) (*Environment, error) {
	var resp schema.EnvironmentGetResponse
	req, err := c.client.NewRequest(ctx, "GET", opts.URL(), nil)
	if err != nil {
		return nil, err
	}

	err = c.client.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	return EnvironmentFromSchema(&resp.Environment), err
}

type EnvironmentUpdateOpts struct {
	EnvironmentGetOpts
	EnvironmentBase
}

func (c *EnvironmentClient) Update(ctx context.Context, opts *EnvironmentUpdateOpts) (*Environment, error) {
	body := &schema.EnvironmentUpdateRequest{}
	opts.EnvironmentBase.setRequest(&body.EnvironmentBase)

	req, err := c.client.NewRequest(ctx, "PUT", opts.URL(), &body)
	if err != nil {
		return nil, err
	}

	var resp schema.EnvironmentUpdateResponse
	err = c.client.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	return EnvironmentFromSchema(&resp.Environment), err
}

type EnvironmentDeleteOpts struct {
	EnvironmentGetOpts
}

func (c *EnvironmentClient) Delete(ctx context.Context, opts *EnvironmentDeleteOpts) error {
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
