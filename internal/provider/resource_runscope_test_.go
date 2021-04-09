// Note this source file ends in an '_'; otherwise the compiler
// will treat is as a test file.

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/terraform-providers/terraform-provider-runscope/internal/runscope"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRunscopeTest() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTestCreate,
		ReadContext:   resourceTestRead,
		UpdateContext: resourceTestUpdate,
		DeleteContext: resourceTestDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.SplitN(d.Id(), "/", 2)
				if len(parts) < 2 {
					return nil, fmt.Errorf("test ID for import should be in format bucket_id/test_id")
				}

				d.Set("bucket_id", parts[0])
				d.SetId(parts[1])

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"bucket_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"default_environment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTestCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	opts := runscope.TestCreateOpts{}
	opts.BucketId = d.Get("bucket_id").(string)
	opts.Name = d.Get("name").(string)
	opts.Description = d.Get("description").(string)

	test, err := client.Test.Create(ctx, opts)
	if err != nil {
		return diag.Errorf("Failed to create test: %s", err)
	}

	d.SetId(test.Id)

	return resourceTestRead(ctx, d, meta)
}

func resourceTestRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	opts := runscope.TestGetOpts{
		BucketId: d.Get("bucket_id").(string),
		Id:       d.Id(),
	}

	test, err := client.Test.Get(ctx, opts)
	if err != nil {
		if isNotFound(err) {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Couldn't read test: %s", err)
	}

	d.Set("name", test.Name)
	d.Set("description", test.Description)
	d.Set("default_environment_id", test.DefaultEnvironmentId)
	return nil
}

func resourceTestUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	opts := runscope.TestUpdateOpts{}
	opts.Id = d.Id()
	opts.BucketId = d.Get("bucket_id").(string)
	opts.Name = d.Get("name").(string)
	opts.Description = d.Get("description").(string)
	opts.DefaultEnvironmentId = d.Get("default_environment_id").(string)

	_, err := client.Test.Update(ctx, opts)
	if err != nil {
		return diag.Errorf("Error updating test: %s", err)
	}

	return resourceTestRead(ctx, d, meta)
}

func resourceTestDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	opts := runscope.TestDeleteOpts{
		Id:       d.Id(),
		BucketId: d.Get("bucket_id").(string),
	}

	if err := client.Test.Delete(ctx, opts); err != nil {
		return diag.Errorf("Error deleting test: %s", err)
	}

	return nil
}
