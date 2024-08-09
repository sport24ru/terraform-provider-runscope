package provider

import (
	"context"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRunscopeIntegration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRunscopeIntegrationRead,

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
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRunscopeIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	searchType := d.Get("type").(string)
	filters, filtersOk := d.GetOk("filter")

	integrations, err := client.Integration.List(ctx, &runscope.IntegrationListOpts{TeamId: d.Get("team_uuid").(string)})
	if err != nil {
		return diag.FromErr(err)
	}

	found := &runscope.Integration{}
	for _, integration := range integrations {
		if integration.Type == searchType {
			if filtersOk {
				if !integrationFiltersTest(integration, filters.(*schema.Set)) {
					continue
				}
			}
			found = integration
			break
		}
	}

	if found == nil {
		return diag.Errorf("Unable to locate any integrations with the type: %s", searchType)
	}

	d.SetId(found.UUID)
	d.Set("type", found.Type)
	d.Set("description", found.Description)

	return nil
}

func integrationFiltersTest(integration *runscope.Integration, filters *schema.Set) bool {
	for _, v := range filters.List() {
		m := v.(map[string]interface{})
		passed := false

		for _, e := range m["values"].(*schema.Set).List() {
			switch m["name"].(string) {
			case "id":
				if integration.UUID == e {
					passed = true
				}
			case "type":
				if integration.Type == e {
					passed = true
				}
			default:
				if integration.Description == e {
					passed = true
				}
			}
		}

		if passed {
			continue
		} else {
			return false
		}

	}
	return true
}
