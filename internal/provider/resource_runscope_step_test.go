package provider

import (
	"context"
	"fmt"
	"github.com/terraform-providers/terraform-provider-runscope/internal/runscope"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccStep_basic(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckStepDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testRunscopeStepConfigA, bucketName, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStepMainPageExists("runscope_step.main_page"),
					resource.TestCheckResourceAttr(
						"runscope_step.main_page", "url", "http://example.com"),
					resource.TestCheckResourceAttr("runscope_step.main_page", "variable.#", "2"),
					resource.TestCheckResourceAttr("runscope_step.main_page", "assertion.#", "2"),
					resource.TestCheckResourceAttr("runscope_step.main_page", "header.#", "3"),
				),
			},
		},
	})
}

func TestAccStep_multiple_steps(t *testing.T) {
	teamID := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckStepDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testRunscopeStepConfigMultipleSteps, bucketName, teamID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStepExists("runscope_step.step_a"),
					testAccCheckStepOrder("runscope_test.test_a", "runscope_step.step_a", "runscope_step.step_b"),
					resource.TestCheckResourceAttr(
						"runscope_step.step_a", "url", "http://step_a.com"),
					resource.TestCheckResourceAttr(
						"runscope_step.step_b", "url", "http://step_b.com")),
			},
			{
				ResourceName:      "runscope_step.step_a",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["runscope_step.step_a"]
					if !ok {
						return "", fmt.Errorf("not found runscope_step.step_a")
					}
					return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["bucket_id"], rs.Primary.Attributes["test_id"], rs.Primary.ID), nil
				},
			},
			{
				ResourceName:      "runscope_step.step_b",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["runscope_step.step_b"]
					if !ok {
						return "", fmt.Errorf("not found runscope_step.step_b")
					}
					return fmt.Sprintf("%s/%s#%d", rs.Primary.Attributes["bucket_id"], rs.Primary.Attributes["test_id"], 2), nil
				},
			},
		},
	})
}

func testAccCheckStepDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*providerConfig).client
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "runscope_step" {
			continue
		}

		opts := &runscope.StepDeleteOpts{}
		opts.BucketId = rs.Primary.Attributes["bucket_id"]
		opts.TestId = rs.Primary.Attributes["test_id"]

		err := client.Step.Delete(ctx, opts)
		if err == nil {
			return fmt.Errorf("Record %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckStepMainPageExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*providerConfig).client

		opts := &runscope.StepGetOpts{}
		opts.Id = rs.Primary.ID
		opts.TestId = rs.Primary.Attributes["test_id"]
		opts.BucketId = rs.Primary.Attributes["bucket_id"]

		step, err := client.Step.Get(ctx, opts)
		if err != nil {
			return err
		}

		if step.Id != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		if len(step.Variables) != 2 {
			return fmt.Errorf("Expected %d variables, actual %d", 2, len(step.Variables))
		}

		variable := step.Variables[1]
		if variable.Name != "httpContentEncoding" {
			return fmt.Errorf("Expected %s variables, actual %s", "httpContentEncoding", variable.Name)
		}

		if len(step.Assertions) != 2 {
			return fmt.Errorf("Expected %d assertions, actual %d", 2, len(step.Assertions))
		}

		assertion := step.Assertions[1]
		if assertion.Source != "response_json" {
			return fmt.Errorf("Expected assertion source %s, actual %s",
				"response_json", assertion.Source)
		}

		if len(step.Headers) != 2 {
			return fmt.Errorf("Expected %d headers, actual %d", 1, len(step.Headers))
		}

		if header, ok := step.Headers["Accept-Encoding"]; ok {
			if len(header) != 2 {
				return fmt.Errorf("Expected %d values for header %s, actual %d",
					2, "Accept-Encoding", len(header))

			}

			if header[1] != "application/xml" {
				return fmt.Errorf("Expected header value %s, actual %s",
					"application/xml", header[1])
			}
		} else {
			return fmt.Errorf("Expected header %s to exist", "Accept-Encoding")
		}

		if len(step.Scripts) != 2 {
			return fmt.Errorf("Expected %d scripts, actual %d", 2, len(step.Scripts))
		}

		if step.Scripts[1] != "log(\"script 2\");" {
			return fmt.Errorf("Expected %s, actual %s", "log(\"script 2\");", step.Scripts[1])
		}

		if len(step.BeforeScripts) != 1 {
			return fmt.Errorf("Expected %d scripts, actual %d", 1, len(step.BeforeScripts))
		}

		if step.BeforeScripts[0] != "log(\"before script\");" {
			return fmt.Errorf("Expected %s, actual %s", "log(\"before script\");", step.BeforeScripts[0])
		}

		if step.Note != "Testing step, single step test" {
			return fmt.Errorf("Expected note %s, actual note %s", "Testing step, single step test", step.Note)
		}

		return nil
	}
}

func testAccCheckStepExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*providerConfig).client

		opts := &runscope.StepGetOpts{}
		opts.Id = rs.Primary.ID
		opts.TestId = rs.Primary.Attributes["test_id"]
		opts.BucketId = rs.Primary.Attributes["bucket_id"]

		step, err := client.Step.Get(ctx, opts)
		if err != nil {
			return err
		}

		if step.Id != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		return nil
	}
}

func testAccCheckStepOrder(n, s1, s2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*providerConfig).client

		opts := runscope.TestGetOpts{
			BucketId: rs.Primary.Attributes["bucket_id"],
			Id:       rs.Primary.ID,
		}

		test, err := client.Test.Get(ctx, opts)
		if err != nil {
			return err
		}

		if test.Id != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		step1, ok := s.RootModule().Resources[s1]
		if !ok {
			return fmt.Errorf("Step not found: %s", s1)
		}

		step2, ok := s.RootModule().Resources[s2]
		if !ok {
			return fmt.Errorf("Step not found: %s", s2)
		}

		if step1.Primary.ID != test.Steps[0].Id {
			return fmt.Errorf("Steps not in correct order, want %s got %s", step1.Primary.ID, test.Steps[0].Id)
		}

		if step2.Primary.ID != test.Steps[1].Id {
			return fmt.Errorf("Steps not in correct order, want %s got %s", step2.Primary.ID, test.Steps[1].Id)
		}

		return nil
	}
}

const testRunscopeStepConfigA = `
resource "runscope_bucket" "bucket" {
  name      = "%s"
  team_uuid = "%s"
}

resource "runscope_test" "test" {
  bucket_id   = "${runscope_bucket.bucket.id}"
  name        = "runscope test"
  description = "This is a test test..."
}

resource "runscope_step" "main_page" {
  bucket_id      = runscope_bucket.bucket.id
  test_id        = runscope_test.test.id
  step_type      = "request"
  note           = "Testing step, single step test"
  url            = "http://example.com"
  method         = "GET"
  variable {
  	   name     = "httpStatus"
  	   source   = "response_status"
  	}
  	variable {
  	   name     = "httpContentEncoding"
  	   source   = "response_header"
  	   property = "Content-Encoding"
  	}

  assertion {
  	   source     = "response_status"
           comparison = "equal_number"
           value      = "200"
  	}
  	assertion {
  	   source     = "response_json"
           comparison = "equal"
           value      = "c5baeb4a-2379-478a-9cda-1b671de77cf9"
           property   = "data.id"
  	}

  header 	{
  		header = "Accept-Encoding"
  		value  = "application/json"
  	}
  	header {
  		header = "Accept-Encoding"
  		value  = "application/xml"
  	}
  	header {
  		header = "Authorization"
  		value  = "Bearer bb74fe7b-b9f2-48bd-9445-bdc60e1edc6a"
	}


  auth {
	username  = "user"
	auth_type = "basic"
	password  = "password1"
  }

  scripts = [
    "log(\"script 1\");",
    "log(\"script 2\");",
  ]
  before_scripts = [
    "log(\"before script\");",
  ]
}
`

const testRunscopeStepConfigMultipleSteps = `
resource "runscope_bucket" "bucket" {
  name      = "%s"
  team_uuid = "%s"
}

resource "runscope_test" "test_a" {
  bucket_id   = runscope_bucket.bucket.id
  name        = "runscope test a"
  description = "This is a test a"
}

resource "runscope_step" "step_b" {
  bucket_id      = runscope_bucket.bucket.id
  test_id        = runscope_test.test_a.id
  step_type      = "request"
  note           = "Multiple step test, test b"
  url            = "http://step_b.com"
  method         = "GET"
  depends_on     = ["runscope_step.step_a"]
}

resource "runscope_step" "step_a" {
  bucket_id      = runscope_bucket.bucket.id
  test_id        = runscope_test.test_a.id
  step_type      = "request"
  note           = "Multiple step test, test a"
  url            = "http://step_a.com"
  method         = "GET"
}
`
