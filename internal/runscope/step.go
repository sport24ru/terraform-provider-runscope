package runscope

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/terraform-providers/terraform-provider-runscope/internal/runscope/schema"
)

type StepBase struct {
	StepType      string
	Method        string
	StepURL       string
	Variables     []StepVariable
	Assertions    []StepAssertion
	Headers       map[string][]string
	Auth          StepAuth
	Body          string
	Form          map[string][]string
	Scripts       []string
	BeforeScripts []string
	Note          string
	Skipped       bool
	Duration      int
}

func (sb *StepBase) setFromSchema(s *schema.Step) {
	sb.StepType = s.StepType
	sb.Method = s.Method
	sb.StepURL = s.URL
	sb.Variables = make([]StepVariable, len(s.Variables))
	sb.Assertions = make([]StepAssertion, len(s.Assertions))
	sb.Headers = map[string][]string{}
	sb.Auth = StepAuth{
		Username: s.Auth.Username,
		Password: s.Auth.Password,
		AuthType: s.Auth.AuthType,
	}
	sb.Body = s.Body
	sb.Form = map[string][]string{}
	sb.Scripts = make([]string, len(s.Scripts))
	sb.BeforeScripts = make([]string, len(s.BeforeScripts))
	sb.Note = s.Note
	sb.Skipped = s.Skipped
	sb.Duration = s.Duration

	for i, v := range s.Variables {
		sb.Variables[i] = StepVariable{
			Name:     v.Name,
			Property: v.Property,
			Source:   v.Source,
		}
	}
	for i, a := range s.Assertions {
		var value string
		if len(a.Value) > 1 && a.Value[0] == '"' && a.Value[len(a.Value)-1] == '"' {
			value = string(a.Value[1 : len(a.Value)-1])
		} else {
			value = string(a.Value)
		}

		sb.Assertions[i] = StepAssertion{
			Source:     a.Source,
			Property:   a.Property,
			Comparison: a.Comparison,
			Value:      value,
		}
	}
	for header, values := range s.Headers {
		sb.Headers[header] = make([]string, len(values))
		for i, v := range values {
			sb.Headers[header][i] = v
		}
	}
	for name, values := range s.Form {
		sb.Form[name] = make([]string, len(values))
		for i, v := range values {
			sb.Form[name][i] = v
		}
	}
	for i, s := range s.Scripts {
		sb.Scripts[i] = s
	}
	for i, s := range s.BeforeScripts {
		sb.BeforeScripts[i] = s
	}
}

type Step struct {
	StepBase
	Id string
}

type StepVariable struct {
	Name     string
	Property string
	Source   string
}

type StepAssertion struct {
	Source     string
	Property   string
	Comparison string
	Value      string
}

type StepAuth struct {
	Username string
	Password string
	AuthType string
}

func (s StepAuth) Empty() bool {
	return s.Username == "" && s.Password == "" && s.AuthType == ""
}

type StepClient struct {
	client *Client
}

func StepFromSchema(s *schema.Step) *Step {
	step := &Step{
		Id: s.Id,
	}
	step.StepBase.setFromSchema(s)

	return step
}

type StepUriOpts struct {
	BucketId string
	TestId   string
}

func (s StepUriOpts) URL() string {
	return fmt.Sprintf("/buckets/%s/tests/%s/steps", s.BucketId, s.TestId)
}

type StepBaseOpts struct {
	StepType      string
	Method        string
	StepURL       string
	Variables     []StepVariable
	Assertions    []StepAssertion
	Headers       map[string][]string
	Auth          StepAuth
	Body          string
	Form          map[string][]string
	Scripts       []string
	BeforeScripts []string
	Note          string
	Skipped       bool
	Duration      int
}

func (sbo *StepBaseOpts) setRequest(sb *schema.StepBase) {
	sb.StepType = sbo.StepType
	sb.Method = sbo.Method
	sb.URL = sbo.StepURL
	sb.Variables = make([]schema.StepVariable, len(sbo.Variables))
	sb.Assertions = make([]schema.StepAssertion, len(sbo.Assertions))
	sb.Headers = map[string][]string{}
	sb.Auth = schema.StepAuth{
		Username: sbo.Auth.Username,
		Password: sbo.Auth.Password,
		AuthType: sbo.Auth.AuthType,
	}
	sb.Body = sbo.Body
	sb.Form = map[string][]string{}
	sb.Scripts = make([]string, len(sbo.Scripts))
	sb.BeforeScripts = make([]string, len(sbo.BeforeScripts))
	sb.Note = sbo.Note
	sb.Skipped = sbo.Skipped
	sb.Duration = sbo.Duration

	for i, v := range sbo.Variables {
		sb.Variables[i] = schema.StepVariable{
			Name:     v.Name,
			Property: v.Property,
			Source:   v.Source,
		}
	}
	for i, a := range sbo.Assertions {
		sb.Assertions[i] = schema.StepAssertion{
			Source:     a.Source,
			Comparison: a.Comparison,
			Value:      (json.RawMessage)(string('"') + string(a.Value) + string('"')),
			Property:   a.Property,
		}
	}
	for header, values := range sbo.Headers {
		sb.Headers[header] = make([]string, len(values))
		for i, v := range values {
			sb.Headers[header][i] = v
		}
	}
	for name, values := range sbo.Form {
		sb.Form[name] = make([]string, len(values))
		for i, v := range values {
			sb.Form[name][i] = v
		}
	}
	for i, s := range sbo.Scripts {
		sb.Scripts[i] = s
	}
	for i, s := range sbo.BeforeScripts {
		sb.BeforeScripts[i] = s
	}
}

type StepCreateOpts struct {
	StepUriOpts
	StepBaseOpts
}

func (c *StepClient) Create(ctx context.Context, opts *StepCreateOpts) (*Step, error) {
	body := &schema.StepCreateRequest{}
	opts.StepBaseOpts.setRequest(&body.StepBase)

	req, err := c.client.NewRequest(ctx, "POST", opts.URL(), &body)
	if err != nil {
		return nil, err
	}

	var resp schema.StepCreateResponse
	err = c.client.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	return StepFromSchema(&resp.Step[len(resp.Step)-1]), nil
}

type StepGetOpts struct {
	StepUriOpts
	Id string
}

func (opts *StepGetOpts) URL() string {
	return fmt.Sprintf("%s/%s", opts.StepUriOpts.URL(), opts.Id)
}

func (c *StepClient) Get(ctx context.Context, opts *StepGetOpts) (*Step, error) {
	var resp schema.StepGetResponse
	req, err := c.client.NewRequest(ctx, "GET", opts.URL(), nil)
	if err != nil {
		return nil, err
	}

	err = c.client.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	return StepFromSchema(&resp.Step), err
}

type StepUpdateOpts struct {
	StepGetOpts
	StepBaseOpts
}

func (c *StepClient) Update(ctx context.Context, opts *StepUpdateOpts) (*Step, error) {
	body := &schema.StepUpdateRequest{}
	opts.StepBaseOpts.setRequest(&body.StepBase)
	body.Id = opts.Id

	req, err := c.client.NewRequest(ctx, "PUT", opts.URL(), &body)
	if err != nil {
		return nil, err
	}

	var resp schema.StepUpdateResponse
	err = c.client.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	return StepFromSchema(&resp.Step), nil
}

type StepDeleteOpts struct {
	StepGetOpts
}

func (c *StepClient) Delete(ctx context.Context, opts *StepDeleteOpts) error {
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
