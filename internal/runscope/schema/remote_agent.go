package schema

type RemoteAgent struct {
	Id      string `json:"agent_id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type RemoteAgentListResponse struct {
	RemoteAgents []RemoteAgent `json:"data"`
}
