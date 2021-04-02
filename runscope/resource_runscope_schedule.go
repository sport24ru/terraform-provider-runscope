package runscope

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/terraform-providers/terraform-provider-runscope/internal/runscope"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRunscopeSchedule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScheduleCreate,
		ReadContext:   resourceScheduleRead,
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
				ForceNew: true,
			},
			"interval": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"note": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceScheduleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	opts := &runscope.ScheduleCreateOpts{}
	opts.EnvironmentId = d.Get("environment_id").(string)
	expandScheduleBaseOpts(d, &opts.ScheduleBaseOpts)
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
	return nil
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

func expandScheduleBaseOpts(d *schema.ResourceData, opts *runscope.ScheduleBaseOpts) {
	opts.BucketId = d.Get("bucket_id").(string)
	opts.TestId = d.Get("test_id").(string)
}

func expandScheduleGetOpts(d *schema.ResourceData, opts *runscope.ScheduleGetOpts) {
	opts.Id = d.Id()
	expandScheduleBaseOpts(d, &opts.ScheduleBaseOpts)
}
