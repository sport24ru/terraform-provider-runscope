package provider

import (
	"context"
	"time"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRunscopeBuckets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRunscopeBucketsRead,

		Schema: map[string]*schema.Schema{
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
			"keys": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceRunscopeBucketsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	filters, filtersOk := d.GetOk("filter")

	buckets, err := client.Bucket.List(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var keys []string
	for _, bucket := range buckets {
		if filtersOk && !bucketFiltersTest(bucket, filters.(*schema.Set)) {
			continue
		}

		keys = append(keys, bucket.Key)
	}

	d.SetId(time.Now().UTC().String())
	d.Set("keys", keys)

	return nil
}

func bucketFiltersTest(bucket *runscope.Bucket, filters *schema.Set) bool {
	for _, v := range filters.List() {
		m := v.(map[string]interface{})
		passed := false

		for _, e := range m["values"].(*schema.Set).List() {
			switch m["name"].(string) {
			case "key":
				if bucket.Key == e {
					passed = true
				}
			default:
				if bucket.Name == e {
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
