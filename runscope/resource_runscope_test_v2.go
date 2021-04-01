// Note this source file ends in an '_'; otherwise the compiler
// will treat is as a test file.

package runscope

import (
	"fmt"
	"log"
	"strings"

	"github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRunscopeTestV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTestV2Create,
		Read:   resourceTestV2Read,
		Update: resourceTestV2Update,
		Delete: resourceTestV2Delete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
				Required: true,
			},
		},
	}
}

func resourceTestV2Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	defaultEnvironmentId, err := verifiedDefaultEnvironmentId(d, client)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	log.Printf("[INFO] Creating test with name: %s", name)

	test, err := createTestV2FromResourceData(d)
	if err != nil {
		return fmt.Errorf("Failed to create test: %s", err)
	}

	log.Printf("[DEBUG] test create: %#v", test)

	createdTest, err := client.CreateTest(test)
	if err != nil {
		return fmt.Errorf("Failed to create test: %s", err)
	}

	if defaultEnvironmentId != "" {
		createdTest.DefaultEnvironmentID = defaultEnvironmentId
		if _, err := client.UpdateTest(createdTest); err != nil {
			return err
		}
	}

	d.SetId(createdTest.ID)
	log.Printf("[INFO] test ID: %s", d.Id())

	return resourceTestRead(d, meta)
}

func verifiedDefaultEnvironmentId(d *schema.ResourceData, client *runscope.Client) (string, error) {
	bucketId := d.Get("bucket_id").(string)

	id := d.Get("default_environment_id").(string)
	_, err := client.ReadSharedEnvironment(
		&runscope.Environment{ID: id},
		&runscope.Bucket{Key: bucketId},
	)
	return id, err
}

func resourceTestV2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	testFromResource, err := createTestFromResourceData(d)
	if err != nil {
		return fmt.Errorf("Error reading test: %s", err)
	}

	test, err := client.ReadTest(testFromResource)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "403") {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Couldn't find test: %s", err)
	}

	d.Set("name", test.Name)
	d.Set("description", test.Description)
	d.Set("default_environment_id", test.DefaultEnvironmentID)
	return nil
}

func resourceTestV2Update(d *schema.ResourceData, meta interface{}) error {
	testFromResource, err := createTestFromResourceData(d)
	if err != nil {
		return fmt.Errorf("Error updating test: %s", err)
	}

	if d.HasChange("description") || d.HasChange("default_environment_id") {
		client := meta.(*runscope.Client)
		_, err = client.UpdateTest(testFromResource)

		if err != nil {
			return fmt.Errorf("Error updating test: %s", err)
		}
	}

	return nil
}

func resourceTestV2Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	test, err := createTestFromResourceData(d)
	if err != nil {
		return fmt.Errorf("Error deleting test: %s", err)
	}
	log.Printf("[INFO] Deleting test with id: %s name: %s", test.ID, test.Name)

	if err := client.DeleteTest(test); err != nil {
		return fmt.Errorf("Error deleting test: %s", err)
	}

	return nil
}

func createTestV2FromResourceData(d *schema.ResourceData) (*runscope.Test, error) {
	test := runscope.NewTest()
	test.ID = d.Id()
	if attr, ok := d.GetOk("bucket_id"); ok {
		test.Bucket.Key = attr.(string)
	}

	if attr, ok := d.GetOk("name"); ok {
		test.Name = attr.(string)
	}

	if attr, ok := d.GetOk("description"); ok {
		test.Description = attr.(string)
	}

	if attr, ok := d.GetOk("default_environment_id"); ok {
		test.DefaultEnvironmentID = attr.(string)
	}

	return test, nil
}
