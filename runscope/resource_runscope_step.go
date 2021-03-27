package runscope

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRunscopeStep() *schema.Resource {
	return &schema.Resource{
		Create: resourceStepCreate,
		Read:   resourceStepRead,
		Update: resourceStepUpdate,
		Delete: resourceStepDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")

				bucketId := parts[0]
				d.Set("bucket_id", bucketId)

				if len(parts) == 3 {
					d.Set("test_id", parts[1])
					d.SetId(parts[2])
					return []*schema.ResourceData{d}, nil
				}

				if len(parts) != 2 {
					return nil, fmt.Errorf("step ID for import should be in format bucket_id/test_id/step_id " +
						"or bucket_id/test_id#step_position")
				}

				parts = strings.Split(parts[1], "#")
				if len(parts) != 2 {
					return nil, fmt.Errorf("step ID for import should be in format bucket_id/test_id/step_id " +
						"or bucket_id/test_id#step_position")
				}

				stepPos, err := strconv.Atoi(parts[1])
				if err != nil || stepPos < 1 {
					return nil, fmt.Errorf("step_position should be a positive integer number")
				}

				testId := parts[0]
				d.Set("test_id", testId)

				test := &runscope.Test{
					ID: testId,
					Bucket: &runscope.Bucket{
						Key: bucketId,
					},
				}

				client := meta.(*runscope.Client)

				test, err = client.ReadTest(test)
				if err != nil {
					return nil, err
				}

				nSteps := len(test.Steps)
				if nSteps < stepPos {
					return nil, fmt.Errorf("test %s contains only %d steps", testId, nSteps)
				}

				d.SetId(test.Steps[stepPos-1].ID)

				return []*schema.ResourceData{d}, nil
			},
		},
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
			"step_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"method": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			// TODO: rename to "variable" for better UX
			"variables": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"property": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"source": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Optional: true,
			},
			// TODO: rename to "assertion" for better UX
			"assertions": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:     schema.TypeString,
							Required: true,
						},
						"property": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"comparison": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
			},
			// TODO: rename to "header" for better UX
			"headers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"auth": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"auth_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"password": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"body": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scripts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"before_scripts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"note": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceStepCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	step, bucketID, testID, err := createStepFromResourceData(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] step create: %#v", step)

	createdStep, err := client.CreateTestStep(step, bucketID, testID)
	if err != nil {
		return fmt.Errorf("Failed to create step: %s", err)
	}

	d.SetId(createdStep.ID)
	log.Printf("[INFO] step ID: %s", d.Id())

	return resourceStepRead(d, meta)
}

func resourceStepRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	stepFromResource, bucketID, testID, err := createStepFromResourceData(d)
	if err != nil {
		return fmt.Errorf("Failed to read step from resource data: %s", err)
	}

	step, err := client.ReadTestStep(stepFromResource, bucketID, testID)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "403") {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Couldn't find step: %s", err)
	}

	d.Set("bucket_id", bucketID)
	d.Set("test_id", testID)
	d.Set("step_type", step.StepType)
	d.Set("method", step.Method)
	d.Set("url", step.URL)
	d.Set("body", step.Body)
	d.Set("variables", readVariables(step.Variables))
	d.Set("assertions", readAssertions(step.Assertions))
	d.Set("headers", readHeaders(step.Headers))
	d.Set("scripts", step.Scripts)
	d.Set("before_scripts", step.BeforeScripts)
	d.Set("note", step.Note)
	if step.Auth != nil && len(step.Auth) > 0 {
		d.Set("auth", []interface{}{
			map[string]interface{}{
				"username":  step.Auth["username"],
				"auth_type": step.Auth["auth_type"],
				"password":  step.Auth["password"],
			},
		})
	}

	return nil
}

func resourceStepUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)
	stepFromResource, bucketID, testID, err := createStepFromResourceData(d)
	if err != nil {
		return fmt.Errorf("Error updating step: %s", err)
	}

	if d.HasChange("url") ||
		d.HasChange("variables") ||
		d.HasChange("assertions") ||
		d.HasChange("headers") ||
		d.HasChange("body") ||
		d.HasChange("note") {
		client := meta.(*runscope.Client)
		_, err = client.UpdateTestStep(stepFromResource, bucketID, testID)

		if err != nil {
			return fmt.Errorf("Error updating step: %s", err)
		}
	}

	return nil
}

func resourceStepDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	stepFromResource, bucketID, testID, err := createStepFromResourceData(d)
	if err != nil {
		return fmt.Errorf("Failed to read step from resource data: %s", err)
	}

	err = client.DeleteTestStep(stepFromResource, bucketID, testID)
	if err != nil {
		return fmt.Errorf("Error deleting step: %s", err)
	}

	return nil
}

func createStepFromResourceData(d *schema.ResourceData) (*runscope.TestStep, string, string, error) {

	step := runscope.NewTestStep()
	bucketID := d.Get("bucket_id").(string)
	testID := d.Get("test_id").(string)
	step.ID = d.Id()
	step.StepType = d.Get("step_type").(string)
	step.Body = d.Get("body").(string)
	if attr, ok := d.GetOk("method"); ok {
		step.Method = attr.(string)
	}

	if attr, ok := d.GetOk("url"); ok {
		step.URL = attr.(string)
	}

	if attr, ok := d.GetOk("variables"); ok {
		variables := []*runscope.Variable{}
		items := attr.(*schema.Set)
		for _, x := range items.List() {
			item := x.(map[string]interface{})
			variable := runscope.Variable{
				Name:     item["name"].(string),
				Property: item["property"].(string),
				Source:   item["source"].(string),
			}

			variables = append(variables, &variable)
		}
		step.Variables = variables
	}

	if v, _ := d.GetOk("auth"); v != nil {
		authSet := v.(*schema.Set).List()
		if len(authSet) == 1 {
			authMap := authSet[0].(map[string]interface{})
			auth := make(map[string]string)
			for key, value := range authMap {
				auth[key] = value.(string)
			}
			step.Auth = auth
		}
	}

	if attr, ok := d.GetOk("assertions"); ok {
		assertions := []*runscope.Assertion{}
		items := attr.([]interface{})
		for _, x := range items {
			item := x.(map[string]interface{})
			variable := runscope.Assertion{
				Source:     item["source"].(string),
				Property:   item["property"].(string),
				Comparison: item["comparison"].(string),
				Value:      item["value"].(string),
			}

			assertions = append(assertions, &variable)
		}

		step.Assertions = assertions
	}

	if attr, ok := d.GetOk("headers"); ok {
		step.Headers = make(map[string][]string)
		items := attr.(*schema.Set)
		for _, x := range items.List() {
			item := x.(map[string]interface{})
			header := item["header"].(string)
			step.Headers[header] = append(step.Headers[header], item["value"].(string))
		}
	}

	if attr, ok := d.GetOk("scripts"); ok {
		step.Scripts = expandStringList(attr.([]interface{}))
	}

	if attr, ok := d.GetOk("before_scripts"); ok {
		step.BeforeScripts = expandStringList(attr.([]interface{}))
	}

	if attr, ok := d.GetOk("note"); ok {
		step.Note = attr.(string)
	}

	return step, bucketID, testID, nil
}

func readVariables(variables []*runscope.Variable) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(variables))
	for _, integration := range variables {

		item := map[string]interface{}{
			"name":     integration.Name,
			"source":   integration.Source,
			"property": integration.Property,
		}

		result = append(result, item)
	}

	return result
}

func readAssertions(assertions []*runscope.Assertion) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(assertions))
	for _, assertion := range assertions {

		item := map[string]interface{}{
			"source":     assertion.Source,
			"property":   assertion.Property,
			"comparison": assertion.Comparison,
			"value":      assertion.Value,
		}

		result = append(result, item)
	}

	return result
}

func readHeaders(headers map[string][]string) []map[string]interface{} {
	result := make([]map[string]interface{}, len(headers))
	for key, header := range headers {
		result = append(result, map[string]interface{}{
			"header": key,
			"value":  header,
		})
	}

	return result
}
