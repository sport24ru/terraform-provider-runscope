package provider

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRunscopeRemoteAgents_Basic(t *testing.T) {
	teamID := os.Getenv("RUNSCOPE_TEAM_ID")
	checks := []resource.TestCheckFunc{}
	i := 0
	for {
		remoteAgentData, ok := os.LookupEnv(fmt.Sprintf("RUNSCOPE_REMOTE_AGENT_%d", i))
		if !ok {
			break
		}
		remoteAgentProps := strings.Split(remoteAgentData, ":")
		checks = append(checks, resource.TestCheckResourceAttr("data.runscope_remote_agents.all", fmt.Sprintf("remote_agents.%d.id", i), remoteAgentProps[0]))
		checks = append(checks, resource.TestCheckResourceAttr("data.runscope_remote_agents.all", fmt.Sprintf("remote_agents.%d.name", i), remoteAgentProps[1]))
		checks = append(checks, resource.TestCheckResourceAttr("data.runscope_remote_agents.all", fmt.Sprintf("remote_agents.%d.version", i), remoteAgentProps[2]))
		i++
	}
	checks = append(checks, resource.TestCheckResourceAttr("data.runscope_remote_agents.all", "remote_agents.#", fmt.Sprintf("%d", i)))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataRemoteAgentsConfig, teamID),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

const testAccDataRemoteAgentsConfig = `
data "runscope_remote_agents" "all" {
  team_uuid = "%s"
}
`
