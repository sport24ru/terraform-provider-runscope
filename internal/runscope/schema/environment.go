package schema

type EnvironmentBase struct {
	Name                string                   `json:"name"`
	Script              string                   `json:"script"`
	PreserveCookies     bool                     `json:"preserve_cookies"`
	InitialVariables    map[string]string        `json:"initial_variables"`
	Integrations        []EnvironmentIntegration `json:"integrations"`
	Regions             []string                 `json:"regions"`
	RemoteAgents        []EnvironmentRemoteAgent `json:"remote_agents"`
	RetryOnFailure      bool                     `json:"retry_on_failure"`
	StopOnFailure       bool                     `json:"stop_on_failure"`
	VerifySSL           bool                     `json:"verify_ssl"`
	Webhooks            []string                 `json:"webhooks"`
	Emails              Emails                   `json:"emails"`
	ParentEnvironmentId string                   `json:"parent_environment_id,omitempty"`
	ClientCertificate   string                   `json:"client_certificate"`
}

type Environment struct {
	EnvironmentBase
	Id string `json:"id"`
}

type EnvironmentIntegration struct {
	Id              string `json:"id"`
	IntegrationType string `json:"integration_type"`
	Description     string `json:"description"`
}

type EnvironmentRemoteAgent struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

type Emails struct {
	NotifyAll       bool        `json:"notify_all"`
	NotifyOn        string      `json:"notify_on"`
	NotifyThreshold int         `json:"notify_threshold"`
	Recipients      []Recipient `json:"recipients"`
}

type Recipient struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type EnvironmentGetResponse struct {
	Environment `json:"data"`
}

type EnvironmentCreateRequest struct {
	EnvironmentBase
}

type EnvironmentCreateResponse struct {
	Environment `json:"data"`
}

type EnvironmentUpdateRequest struct {
	EnvironmentBase
}

type EnvironmentUpdateResponse struct {
	Environment `json:"data"`
}
