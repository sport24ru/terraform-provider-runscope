package provider

import (
	"context"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRunscopeBucket() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRunscopeBucketRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"team_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"verify_ssl": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"trigger_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRunscopeBucketRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	opts := &runscope.BucketGetOpts{Key: d.Get("key").(string)}
	bucket, err := client.Bucket.Get(ctx, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(bucket.Key)
	d.Set("name", bucket.Name)
	d.Set("team_uuid", bucket.Team.UUID)
	d.Set("auth_token", bucket.AuthToken)
	d.Set("default", bucket.Default)
	d.Set("verify_ssl", bucket.VerifySSL)
	d.Set("trigger_url", bucket.TriggerURL)

	return nil
}
