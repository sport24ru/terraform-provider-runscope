package provider

import (
	"context"
	"time"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRunscopeIntegrations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRunscopeIntegrationsRead,

		Schema: map[string]*schema.Schema{
			"team_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceRunscopeIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	filters, filtersOk := d.GetOk("filter")

	integrations, err := client.Integration.List(ctx, &runscope.IntegrationListOpts{TeamId: d.Get("team_uuid").(string)})
	if err != nil {
		return diag.FromErr(err)
	}

	var ids []string
	for _, integration := range integrations {
		if filtersOk {
			if !integrationFiltersTest(integration, filters.(*schema.Set)) {
				continue
			}
		}

		ids = append(ids, integration.UUID)
	}

	d.SetId(time.Now().UTC().String())
	d.Set("ids", ids)

	return nil
}
