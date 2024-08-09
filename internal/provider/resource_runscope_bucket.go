package provider

import (
	"context"
	"fmt"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRunscopeBucket() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBucketCreate,
		ReadContext:   resourceBucketRead,
		DeleteContext: resourceBucketDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBucketImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"team_uuid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

func resourceBucketCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	opts := &runscope.BucketCreateOpts{
		Name:     d.Get("name").(string),
		TeamUUID: d.Get("team_uuid").(string),
	}

	bucket, err := client.Bucket.Create(ctx, opts)
	if err != nil {
		return diag.Errorf("Failed to create bucket: %s", err)
	}

	d.SetId(bucket.Key)

	return resourceBucketRead(ctx, d, meta)
}

func resourceBucketRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	opts := &runscope.BucketGetOpts{Key: d.Id()}
	bucket, err := client.Bucket.Get(ctx, opts)
	if err != nil {
		if isNotFound(err) {
			d.SetId("")
			return nil
		}

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

func resourceBucketImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	key := d.Id()

	diags := resourceBucketRead(ctx, d, meta)
	if diags.HasError() {
		return nil, diags[0].Validate()
	}

	if d.Id() == "" {
		return nil, fmt.Errorf("Couldn't find bucket: %s", key)
	}

	results := []*schema.ResourceData{d}

	return results, nil
}

func resourceBucketDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	opts := &runscope.BucketDeleteOpts{}
	opts.Key = d.Id()

	if err := client.Bucket.Delete(ctx, opts); err != nil {
		return diag.Errorf("Error deleting bucket: %s", err)
	}

	return nil
}
