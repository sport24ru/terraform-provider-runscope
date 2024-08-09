package runscope

import (
	"context"
	"fmt"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope/schema"
)

type RemoteAgent struct {
	Id      string
	Name    string
	Version string
}

type RemoteAgentClient struct {
	client *Client
}

func RemoteAgentFromSchema(s schema.RemoteAgent) *RemoteAgent {
	return &RemoteAgent{
		Id:      s.Id,
		Name:    s.Name,
		Version: s.Version,
	}
}

type RemoteAgentListOpts struct {
	TeamUUID string
}

// URL returns an URL of agents list request
//
// See https://api.blazemeter.com/api-monitoring/#team-agents-list
func (opts *RemoteAgentListOpts) URL() string {
	return fmt.Sprintf("/teams/%s/agents", opts.TeamUUID)
}

// List returns a list of the team’s currently connected agents.
//
// See https://api.blazemeter.com/api-monitoring/#team-agents-list
func (c *RemoteAgentClient) List(ctx context.Context, opts *RemoteAgentListOpts) ([]*RemoteAgent, error) {
	req, err := c.client.NewRequest(ctx, "GET", opts.URL(), nil)
	if err != nil {
		return nil, err
	}

	var resp schema.RemoteAgentListResponse
	err = c.client.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	remoteAgents := make([]*RemoteAgent, len(resp.RemoteAgents), len(resp.RemoteAgents))
	for i, remoteAgent := range resp.RemoteAgents {
		remoteAgents[i] = RemoteAgentFromSchema(remoteAgent)
	}

	return remoteAgents, nil
}
