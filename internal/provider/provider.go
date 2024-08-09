package provider

import (
	"context"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RUNSCOPE_ACCESS_TOKEN", nil),
				Description: "A runscope access token.",
			},
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RUNSCOPE_API_URL", nil),
				Description: "A runscope api url i.e. https://api.runscope.com.",
				Default:     "https://api.runscope.com",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"runscope_integration":   dataSourceRunscopeIntegration(),
			"runscope_integrations":  dataSourceRunscopeIntegrations(),
			"runscope_bucket":        dataSourceRunscopeBucket(),
			"runscope_buckets":       dataSourceRunscopeBuckets(),
			"runscope_remote_agents": dataSourceRunscopeRemoteAgents(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"runscope_bucket":      resourceRunscopeBucket(),
			"runscope_test":        resourceRunscopeTest(),
			"runscope_environment": resourceRunscopeEnvironment(),
			"runscope_schedule":    resourceRunscopeSchedule(),
			"runscope_step":        resourceRunscopeStep(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

type providerConfig struct {
	client *runscope.Client
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("access_token").(string)
	endpoint := d.Get("api_url").(string)

	client := runscope.NewClient(runscope.WithToken(token), runscope.WithEndpoint(endpoint))

	return &providerConfig{
		client: client,
	}, nil
}

func isNotFound(err error) bool {
	if runscopeErr, ok := err.(runscope.Error); ok {
		return runscopeErr.Status() == 404
	}
	return false
}
