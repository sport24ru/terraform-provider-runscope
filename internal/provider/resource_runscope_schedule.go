package provider

import (
	"context"
	"strings"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRunscopeSchedule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScheduleCreate,
		ReadContext:   resourceScheduleRead,
		UpdateContext: resourceScheduleUpdate,
		DeleteContext: resourceScheduleDelete,

		Schema: map[string]*schema.Schema{
			"bucket_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"test_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"interval": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"1m", "5m", "15m", "30m", "1h", "6h", "1d"}, false),
			},
			"note": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"exported_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceScheduleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	opts := &runscope.ScheduleCreateOpts{}
	opts.EnvironmentId = d.Get("environment_id").(string)
	expandScheduleBaseOpts(d, &opts.ScheduleURLOpts)
	if v, ok := d.GetOk("interval"); ok {
		opts.Interval = v.(string)
	}
	if v, ok := d.GetOk("note"); ok {
		opts.Note = v.(string)
	}

	schedule, err := client.Schedule.Create(ctx, opts)
	if err != nil {
		if isNotFound(err) {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Couldn't create schedule: %s", err)
	}

	d.SetId(schedule.Id)

	return resourceScheduleRead(ctx, d, meta)
}

func resourceScheduleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	opts := &runscope.ScheduleGetOpts{}
	expandScheduleGetOpts(d, opts)

	schedule, err := client.Schedule.Get(ctx, opts)
	if err != nil {
		if isNotFound(err) {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Couldn't read schedule: %s", err)
	}

	d.Set("environment_id", schedule.EnvironmentId)
	d.Set("interval", strings.ReplaceAll(schedule.Interval, ".0", ""))
	d.Set("note", schedule.Note)
	d.Set("exported_at", flattenTime(schedule.ExportedAt))
	return nil
}

func resourceScheduleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	opts := runscope.ScheduleUpdateOpts{}
	opts.Id = d.Id()
	opts.BucketId = d.Get("bucket_id").(string)
	opts.TestId = d.Get("test_id").(string)
	opts.EnvironmentId = d.Get("environment_id").(string)
	opts.Interval = d.Get("interval").(string)
	opts.Note = d.Get("note").(string)

	_, err := client.Schedule.Update(ctx, &opts)
	if err != nil {
		return diag.Errorf("Error updating schedule: %s", err)
	}

	return resourceScheduleRead(ctx, d, meta)
}

func resourceScheduleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	opts := &runscope.ScheduleDeleteOpts{}
	expandScheduleGetOpts(d, &opts.ScheduleGetOpts)

	if err := client.Schedule.Delete(ctx, opts); err != nil {
		return diag.Errorf("Error deleting test: %s", err)
	}

	return nil
}

func expandScheduleBaseOpts(d *schema.ResourceData, opts *runscope.ScheduleURLOpts) {
	opts.BucketId = d.Get("bucket_id").(string)
	opts.TestId = d.Get("test_id").(string)
}

func expandScheduleGetOpts(d *schema.ResourceData, opts *runscope.ScheduleGetOpts) {
	opts.Id = d.Id()
	expandScheduleBaseOpts(d, &opts.ScheduleURLOpts)
}
