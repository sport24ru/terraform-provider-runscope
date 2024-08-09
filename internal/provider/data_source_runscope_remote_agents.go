package provider

import (
	"context"
	"time"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRunscopeRemoteAgents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRunscopeRemoteAgentsRead,

		Schema: map[string]*schema.Schema{
			"team_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remote_agents": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRunscopeRemoteAgentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	remoteAgents, err := client.RemoteAgent.List(ctx, &runscope.RemoteAgentListOpts{TeamUUID: d.Get("team_uuid").(string)})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("remote_agents", flattenRemoteAgents(remoteAgents)); err != nil {
		return diag.Errorf("error setting remote_agents for data.runscope_remote_agents %s: %s", d.Id(), err)
	}

	return nil
}
