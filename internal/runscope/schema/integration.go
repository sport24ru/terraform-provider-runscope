package schema

type Integration struct {
	UUID        string `json:"uuid"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type IntegrationListResponse struct {
	Integrations []Integration `json:"data"`
}
