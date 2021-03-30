package runscope

import (
	"fmt"
	"log"
	"strings"

	"github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceRunscopeEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnvironmentCreate,
		Read:   resourceEnvironmentRead,
		Update: resourceEnvironmentUpdate,
		Delete: resourceEnvironmentDelete,

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
			"test_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"script": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"preserve_cookies": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"initial_variables": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: false,
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
		},
	}
}

func resourceEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	name := d.Get("name").(string)
	log.Printf("[INFO] Creating environment with name: %s", name)

	environment := expandEnvironment(d)
	log.Printf("[DEBUG] environment create: %#v", environment)

	var createdEnvironment *runscope.Environment
	bucketID := d.Get("bucket_id").(string)

	var err error
	if testID, ok := d.GetOk("test_id"); ok {
		createdEnvironment, err = client.CreateTestEnvironment(environment,
			&runscope.Test{ID: testID.(string), Bucket: &runscope.Bucket{Key: bucketID}})
	} else {
		createdEnvironment, err = client.CreateSharedEnvironment(environment,
			&runscope.Bucket{Key: bucketID})
	}
	if err != nil {
		return fmt.Errorf("Failed to create environment: %s", err)
	}

	d.SetId(createdEnvironment.ID)
	log.Printf("[INFO] environment ID: %s", d.Id())

	return resourceEnvironmentRead(d, meta)
}

func resourceEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	environmentFromResource := expandEnvironment(d)

	var environment *runscope.Environment
	var err error
	bucketID := d.Get("bucket_id").(string)
	if testID, ok := d.GetOk("test_id"); ok {
		environment, err = client.ReadTestEnvironment(
			environmentFromResource, &runscope.Test{ID: testID.(string), Bucket: &runscope.Bucket{Key: bucketID}})
	} else {
		environment, err = client.ReadSharedEnvironment(
			environmentFromResource, &runscope.Bucket{Key: bucketID})
	}

	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "403") {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Couldn't find environment: %s", err)
	}

	d.Set("bucket_id", bucketID)
	d.Set("test_id", d.Get("test_id").(string))
	d.Set("name", environment.Name)
	d.Set("script", environment.Script)
	d.Set("preserve_cookies", environment.PreserveCookies)
	d.Set("initial_variables", environment.InitialVariables)
	d.Set("integrations", flattenIntegrations(environment.Integrations))
	d.Set("retry_on_failure", environment.RetryOnFailure)
	d.Set("verify_ssl", environment.VerifySsl)
	d.Set("webhooks", environment.WebHooks)

	d.Set("email", flattenEmailSettings(environment.EmailSettings))
	return nil
}

func resourceEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	environment := expandEnvironment(d)

	client := meta.(*runscope.Client)
	bucketID := d.Get("bucket_id").(string)
	var err error
	if testID, ok := d.GetOk("test_id"); ok {
		_, err = client.UpdateTestEnvironment(
			environment, &runscope.Test{ID: testID.(string), Bucket: &runscope.Bucket{Key: bucketID}})
	} else {
		_, err = client.UpdateSharedEnvironment(
			environment, &runscope.Bucket{Key: bucketID})
	}
	if err != nil {
		return fmt.Errorf("Error updating environment: %s", err)
	}

	return resourceEnvironmentRead(d, meta)
}

func resourceEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	environmentFromResource := expandEnvironment(d)

	bucketID := d.Get("bucket_id").(string)
	var err error

	if testID, ok := d.GetOk("test_id"); ok {
		log.Printf("[INFO] Deleting test environment with id: %s name: %s, from test %s",
			environmentFromResource.ID, environmentFromResource.Name, testID.(string))
		err = client.DeleteEnvironment(
			environmentFromResource, &runscope.Bucket{Key: bucketID})
	} else {
		log.Printf("[INFO] Deleting shared environment with id: %s name: %s",
			environmentFromResource.ID, environmentFromResource.Name)
		err = client.DeleteEnvironment(
			environmentFromResource, &runscope.Bucket{Key: bucketID})
	}

	if err != nil {
		return fmt.Errorf("Error deleting environment: %s", err)
	}

	return nil
}

func expandEnvironment(d *schema.ResourceData) *runscope.Environment {
	environment := runscope.NewEnvironment()

	environment.ID = d.Id()
	environment.TestID = d.Get("test_id").(string)
	environment.Name = d.Get("name").(string)

	if attr, ok := d.GetOk("script"); ok {
		environment.Script = attr.(string)
	}

	if attr, ok := d.GetOk("preserve_cookies"); ok {
		environment.PreserveCookies = attr.(bool)
	}

	if attr, ok := d.GetOk("initial_variables"); ok {
		variablesRaw := attr.(map[string]interface{})
		variables := map[string]string{}
		for k, v := range variablesRaw {
			variables[k] = v.(string)
		}

		environment.InitialVariables = variables
	}

	if attr, ok := d.GetOk("integrations"); ok {
		integrations := []*runscope.EnvironmentIntegration{}
		items := attr.(*schema.Set)
		for _, item := range items.List() {
			integration := runscope.EnvironmentIntegration{
				ID: item.(string),
			}

			integrations = append(integrations, &integration)
		}

		environment.Integrations = integrations
	}

	if attr, ok := d.GetOk("regions"); ok {
		regions := []string{}
		items := attr.(*schema.Set)
		for _, x := range items.List() {
			item := x.(string)
			regions = append(regions, item)
		}

		environment.Regions = regions
	}

	if attr, ok := d.GetOk("remote_agent"); ok {
		remoteAgents := []*runscope.LocalMachine{}
		items := attr.(*schema.Set)
		for _, x := range items.List() {
			item := x.(map[string]interface{})
			remoteAgent := runscope.LocalMachine{
				Name: item["name"].(string),
				UUID: item["uuid"].(string),
			}

			remoteAgents = append(remoteAgents, &remoteAgent)
		}
		environment.RemoteAgents = remoteAgents
	}

	if attr, ok := d.GetOk("retry_on_failure"); ok {
		environment.RetryOnFailure = attr.(bool)
	}

	environment.VerifySsl = d.Get("verify_ssl").(bool)

	if attr, ok := d.GetOk("webhooks"); ok {
		webhooks := []string{}
		items := attr.(*schema.Set)
		for _, x := range items.List() {
			item := x.(string)
			webhooks = append(webhooks, item)
		}

		environment.WebHooks = webhooks
	}

	environment.EmailSettings = expandEmailSettings(d.Get("email"), false)

	return environment
}

func expandEmailSettings(v interface{}, emptyIsNil bool) *runscope.EmailSettings {
	es := &runscope.EmailSettings{}
	if v == nil {
		if emptyIsNil {
			return nil
		}
		return es
	}

	list := v.([]interface{})
	if len(list) < 1 {
		if emptyIsNil {
			return nil
		}
		return es
	}

	m := list[0].(map[string]interface{})
	es.NotifyAll = m["notify_all"].(bool)
	es.NotifyOn = m["notify_on"].(string)
	es.NotifyThreshold = m["notify_threshold"].(int)
	es.Recipients = expandRecipients(m["recipient"])

	return es
}

func expandRecipients(v interface{}) []*runscope.Contact {
	contacts := make([]*runscope.Contact, 0)
	for _, val := range v.(*schema.Set).List() {
		contacts = append(contacts, &runscope.Contact{
			ID: val.(map[string]interface{})["id"].(string),
		})
	}
	return contacts
}

func flattenIntegrations(integrations []*runscope.EnvironmentIntegration) []interface{} {
	result := []interface{}{}

	for _, integration := range integrations {
		result = append(result, integration.ID)
	}

	return result
}

func flattenEmailSettings(emailSettings *runscope.EmailSettings) []interface{} {
	if isDefaultEmailSettings(emailSettings) {
		return nil
	}

	item := map[string]interface{}{
		"notify_all":       emailSettings.NotifyAll,
		"notify_on":        emailSettings.NotifyOn,
		"notify_threshold": emailSettings.NotifyThreshold,
	}

	if len(emailSettings.Recipients) > 0 {
		resultRecipients := []interface{}{}

		for _, recipient := range emailSettings.Recipients {
			item := map[string]interface{}{
				"name":  recipient.Name,
				"email": recipient.Email,
				"id":    recipient.ID,
			}
			resultRecipients = append(resultRecipients, item)
		}

		item["recipient"] = resultRecipients
	}

	return []interface{}{item}
}

func isDefaultEmailSettings(e *runscope.EmailSettings) bool {
	return e.NotifyAll == false && e.NotifyOn == "" && e.NotifyThreshold == 0 && len(e.Recipients) == 0
}

func recipientsHash(v interface{}) int {
	m := v.(map[string]interface{})
	return schema.HashString(m["id"].(string))
}
