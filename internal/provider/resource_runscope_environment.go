package provider

import (
	"context"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceRunscopeEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,

		Schema: map[string]*schema.Schema{
			"bucket_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"test_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"script": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"preserve_cookies": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"initial_variables": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"integrations": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"regions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"remote_agent": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Optional: true,
			},
			"retry_on_failure": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"stop_on_failure": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"verify_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"webhooks": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"email": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notify_all": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"notify_on": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"all", "failures", "threshold", "switch",
							}, false),
						},
						"notify_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"recipient": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"email": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							Set: recipientsHash,
						},
					},
				},
				Optional: true,
			},
			"parent_environment_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	if err := validateEnvironmentSchema(d); err != nil {
		return err
	}

	opts := runscope.EnvironmentCreateOpts{}
	expandEnvironmentUriOpts(d, &opts.EnvironmentUriOpts)
	expandEnvironmentBase(d, &opts.EnvironmentBase)
	if env, err := client.Environment.Create(ctx, &opts); err != nil {
		return diag.Errorf("Couldn't create environment: %s", err)
	} else {
		d.SetId(env.Id)
	}

	return resourceEnvironmentRead(ctx, d, meta)
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	opts := runscope.EnvironmentGetOpts{
		Id: d.Id(),
	}
	opts.BucketId = d.Get("bucket_id").(string)
	if v, ok := d.GetOk("test_id"); ok {
		opts.TestId = v.(string)
	}

	env, err := client.Environment.Get(ctx, &opts)
	if err != nil {
		if isNotFound(err) {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Couldn't read environment: %s", err)
	}

	d.Set("bucket_id", opts.BucketId)
	d.Set("test_id", opts.TestId)
	d.Set("name", env.Name)
	d.Set("script", env.Script)
	d.Set("preserve_cookies", env.PreserveCookies)
	d.Set("initial_variables", env.InitialVariables)
	d.Set("integrations", env.Integrations)
	d.Set("retry_on_failure", env.RetryOnFailure)
	d.Set("stop_on_failure", env.StopOnFailure)
	d.Set("verify_ssl", env.VerifySSL)
	d.Set("webhooks", env.Webhooks)
	if !env.Emails.IsDefault() {
		d.Set("email", flattenEmails(env.Emails))
	}
	d.Set("parent_environment_id", env.ParentEnvironmentId)
	d.Set("client_certificate", env.ClientCertificate)

	return nil
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	if err := validateEnvironmentSchema(d); err != nil {
		return err
	}

	opts := runscope.EnvironmentUpdateOpts{}
	expandEnvironmentGetOpts(d, &opts.EnvironmentGetOpts)
	expandEnvironmentBase(d, &opts.EnvironmentBase)

	if _, err := client.Environment.Update(ctx, &opts); err != nil {
		return diag.Errorf("Couldn't update environment: %s", err)
	}

	return resourceEnvironmentRead(ctx, d, meta)
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerConfig).client

	opts := runscope.EnvironmentDeleteOpts{}
	expandEnvironmentGetOpts(d, &opts.EnvironmentGetOpts)

	if err := client.Environment.Delete(ctx, &opts); err != nil {
		return diag.Errorf("Error deleting environment: %s", err)
	}

	return nil
}

func flattenEmails(e runscope.Emails) []interface{} {
	return []interface{}{map[string]interface{}{
		"notify_all":       e.NotifyAll,
		"notify_on":        e.NotifyOn,
		"notify_threshold": e.NotifyThreshold,
		"recipient":        flattenRecipients(e.Recipients),
	}}
}

func flattenRecipients(re []runscope.Recipient) []map[string]interface{} {
	recipients := []map[string]interface{}{}
	for _, rec := range re {
		recipients = append(recipients, map[string]interface{}{
			"id":    rec.Id,
			"name":  rec.Name,
			"email": rec.Email,
		})
	}
	return recipients
}

func recipientsHash(v interface{}) int {
	m := v.(map[string]interface{})
	return schema.HashString(m["id"].(string))
}

func expandEnvironmentUriOpts(d *schema.ResourceData, opts *runscope.EnvironmentUriOpts) {
	opts.BucketId = d.Get("bucket_id").(string)
	if v, ok := d.GetOk("test_id"); ok {
		opts.TestId = v.(string)
	}
}

func expandEnvironmentGetOpts(d *schema.ResourceData, opts *runscope.EnvironmentGetOpts) {
	opts.Id = d.Id()
	expandEnvironmentUriOpts(d, &opts.EnvironmentUriOpts)
}

func expandEnvironmentBase(d *schema.ResourceData, opts *runscope.EnvironmentBase) {
	opts.Name = d.Get("name").(string)
	opts.VerifySSL = d.Get("verify_ssl").(bool)
	if v, ok := d.GetOk("script"); ok {
		opts.Script = v.(string)
	}
	if v, ok := d.GetOk("preserve_cookies"); ok {
		opts.PreserveCookies = v.(bool)
	}
	if v, ok := d.GetOk("initial_variables"); ok {
		opts.InitialVariables = map[string]string{}
		for key, value := range v.(map[string]interface{}) {
			opts.InitialVariables[key] = value.(string)
		}
	}
	if v, ok := d.GetOk("integrations"); ok {
		for _, id := range v.(*schema.Set).List() {
			opts.Integrations = append(opts.Integrations, id.(string))
		}
	}
	if v, ok := d.GetOk("regions"); ok {
		for _, region := range v.(*schema.Set).List() {
			opts.Regions = append(opts.Regions, region.(string))
		}
	}
	if v, ok := d.GetOk("remote_agent"); ok {
		for _, ra := range v.(*schema.Set).List() {
			raa := ra.(map[string]interface{})
			opts.RemoteAgents = append(opts.RemoteAgents, runscope.EnvironmentRemoteAgent{
				Name: raa["name"].(string),
				UUID: raa["uuid"].(string),
			})
		}
	}
	if v, ok := d.GetOk("retry_on_failure"); ok {
		opts.RetryOnFailure = v.(bool)
	}
	if v, ok := d.GetOk("stop_on_failure"); ok {
		opts.StopOnFailure = v.(bool)
	}
	if v, ok := d.GetOk("webhooks"); ok {
		for _, w := range v.(*schema.Set).List() {
			opts.Webhooks = append(opts.Webhooks, w.(string))
		}
	}
	if v, ok := d.GetOk("email"); ok {
		emails := v.([]interface{})
		if len(emails) > 0 {
			ee := emails[0].(map[string]interface{})
			opts.Emails.NotifyOn = ee["notify_on"].(string)
			opts.Emails.NotifyAll = ee["notify_all"].(bool)
			opts.Emails.NotifyThreshold = ee["notify_threshold"].(int)
			recipients := ee["recipient"].(*schema.Set).List()
			if len(recipients) > 0 {
				opts.Emails.Recipients = make([]runscope.Recipient, len(recipients))
				for i, recipient := range recipients {
					r := recipient.(map[string]interface{})
					opts.Emails.Recipients[i].Id = r["id"].(string)
					opts.Emails.Recipients[i].Name = r["name"].(string)
					opts.Emails.Recipients[i].Email = r["email"].(string)
				}
			}
		}
	}
	if v, ok := d.GetOk("parent_environment_id"); ok {
		opts.ParentEnvironmentId = v.(string)
	}
	if v, ok := d.GetOk("client_certificate"); ok {
		opts.ClientCertificate = v.(string)
	}
}

func validateEnvironmentSchema(d *schema.ResourceData) diag.Diagnostics {
	if _, hasTestId := d.GetOk("test_id"); hasTestId {
		return nil
	}
	if _, hasParentEnvironmentId := d.GetOk("parent_environment_id"); !hasParentEnvironmentId {
		return nil
	}

	return diag.Errorf("parent_environment_id could be set only if test_id defined")
}
