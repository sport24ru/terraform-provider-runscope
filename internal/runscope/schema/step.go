package schema

import "encoding/json"

type StepBase struct {
	StepType      string              `json:"step_type"`
	Method        string              `json:"method"`
	URL           string              `json:"url"`
	Variables     []StepVariable      `json:"variables"`
	Assertions    []StepAssertion     `json:"assertions"`
	Headers       map[string][]string `json:"headers"`
	Auth          StepAuth            `json:"auth"`
	Body          string              `json:"body"`
	Form          map[string][]string `json:"form"`
	Scripts       []string            `json:"scripts"`
	BeforeScripts []string            `json:"before_scripts"`
	Note          string              `json:"note"`
	Skipped       bool                `json:"skipped"`
	Duration      int                 `json:"duration"`
}

type Step struct {
	StepBase
	Id string `json:"id"`
}

type StepVariable struct {
	Name     string `json:"name"`
	Property string `json:"property"`
	Source   string `json:"source"`
}

type StepAssertion struct {
	Source     string          `json:"source"`
	Property   string          `json:"property"`
	Comparison string          `json:"comparison"`
	Value      json.RawMessage `json:"value"`
}

type StepAuth struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	AuthType string `json:"auth_type,omitempty"`
}

type StepGetResponse struct {
	Step `json:"data"`
}

type StepCreateRequest struct {
	StepBase
}

type StepCreateResponse struct {
	Step []Step `json:"data"`
}

type StepUpdateRequest struct {
	Step
}

type StepUpdateResponse struct {
	Step Step `json:"data"`
}
