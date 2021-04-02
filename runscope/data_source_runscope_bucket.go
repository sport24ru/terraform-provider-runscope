package runscope

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-runscope/internal/runscope"
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

	return nil
}
