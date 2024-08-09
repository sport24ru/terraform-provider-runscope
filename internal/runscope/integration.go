package runscope

import (
	"context"
	"fmt"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope/schema"
)

type Integration struct {
	UUID        string
	Type        string
	Description string
}

type IntegrationClient struct {
	client *Client
}

func IntegrationFromSchema(s schema.Integration) *Integration {
	return &Integration{
		UUID:        s.UUID,
		Type:        s.Type,
		Description: s.Description,
	}
}

type IntegrationListOpts struct {
	TeamId string
}

func (opts *IntegrationListOpts) URL() string {
	return fmt.Sprintf("/teams/%s/integrations", opts.TeamId)
}

func (c *IntegrationClient) List(ctx context.Context, opts *IntegrationListOpts) ([]*Integration, error) {
	req, err := c.client.NewRequest(ctx, "GET", opts.URL(), nil)
	if err != nil {
		return nil, err
	}

	var resp schema.IntegrationListResponse
	err = c.client.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	integrations := make([]*Integration, len(resp.Integrations), len(resp.Integrations))
	for i, integration := range resp.Integrations {
		integrations[i] = IntegrationFromSchema(integration)
	}

	return integrations, nil
}
