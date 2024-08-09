package provider

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccStep_create_default(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckStepDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccStepDefaultConfig, bucketName, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStepExists("runscope_step.step"),
					resource.TestCheckResourceAttr("runscope_step.step", "step_type", "request"),
					resource.TestCheckResourceAttr("runscope_step.step", "skipped", "false"),
					resource.TestCheckResourceAttr("runscope_step.step", "note", ""),
					resource.TestCheckResourceAttr("runscope_step.step", "method", "GET"),
					resource.TestCheckResourceAttr("runscope_step.step", "body", ""),
					resource.TestCheckResourceAttr("runscope_step.step", "form_parameter.#", "0"),
					resource.TestCheckResourceAttr("runscope_step.step", "header.#", "0"),
					resource.TestCheckResourceAttr("runscope_step.step", "url", "https://example.org"),
					resource.TestCheckResourceAttr("runscope_step.step", "auth.#", "0"),
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.#", "0"),
					resource.TestCheckResourceAttr("runscope_step.step", "variable.#", "0"),
					resource.TestCheckResourceAttr("runscope_step.step", "scripts.#", "0"),
					resource.TestCheckResourceAttr("runscope_step.step", "before_scripts.#", "0"),
				),
			},
		},
	})
}

func TestAccStep_create_custom(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckStepDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccStepCustomConfig, bucketName, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStepExists("runscope_step.step"),
					resource.TestCheckResourceAttr("runscope_step.step", "step_type", "request"),
					resource.TestCheckResourceAttr("runscope_step.step", "skipped", "true"),
					resource.TestCheckResourceAttr("runscope_step.step", "note", "Step note"),
					resource.TestCheckResourceAttr("runscope_step.step", "method", "POST"),
					resource.TestCheckResourceAttr("runscope_step.step", "url", "https://example.org"),
					resource.TestCheckResourceAttr("runscope_step.step", "body", "Request body"),
					resource.TestCheckResourceAttr("runscope_step.step", "form_parameter.#", "3"),
					resource.TestCheckResourceAttr("runscope_step.step", "form_parameter.0.name", "a"),
					resource.TestCheckResourceAttr("runscope_step.step", "form_parameter.0.value", "1"),
					resource.TestCheckResourceAttr("runscope_step.step", "form_parameter.1.name", "a"),
					resource.TestCheckResourceAttr("runscope_step.step", "form_parameter.1.value", "2"),
					resource.TestCheckResourceAttr("runscope_step.step", "form_parameter.2.name", "b"),
					resource.TestCheckResourceAttr("runscope_step.step", "form_parameter.2.value", "3"),
					resource.TestCheckResourceAttr("runscope_step.step", "header.#", "3"),
					resource.TestCheckResourceAttr("runscope_step.step", "header.0.header", "Accept-Encoding"),
					resource.TestCheckResourceAttr("runscope_step.step", "header.0.value", "application/json"),
					resource.TestCheckResourceAttr("runscope_step.step", "header.1.header", "Accept-Encoding"),
					resource.TestCheckResourceAttr("runscope_step.step", "header.1.value", "application/xml"),
					resource.TestCheckResourceAttr("runscope_step.step", "header.2.header", "Authorization"),
					resource.TestCheckResourceAttr("runscope_step.step", "header.2.value", "Bearer bb74fe7b-b9f2-48bd-9445-bdc60e1edc6a"),
					resource.TestCheckResourceAttr("runscope_step.step", "auth.#", "1"),
					resource.TestCheckResourceAttr("runscope_step.step", "auth.0.username", "user"),
					resource.TestCheckResourceAttr("runscope_step.step", "auth.0.auth_type", "basic"),
					resource.TestCheckResourceAttr("runscope_step.step", "auth.0.password", "password1"),
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.#", "2"),
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.0.source", "response_status"),
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.0.comparison", "equal_number"),
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.0.value", "200"),
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.1.source", "response_json"),
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.1.comparison", "equal"),
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.1.value", "c5baeb4a-2379-478a-9cda-1b671de77cf9"),
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.1.property", "data.id"),
					resource.TestCheckResourceAttr("runscope_step.step", "variable.#", "2"),
					resource.TestCheckResourceAttr("runscope_step.step", "variable.0.name", "httpContentEncoding"),
					resource.TestCheckResourceAttr("runscope_step.step", "variable.0.source", "response_headers"),
					resource.TestCheckResourceAttr("runscope_step.step", "variable.0.property", "Content-Encoding"),
					resource.TestCheckResourceAttr("runscope_step.step", "variable.1.name", "httpStatus"),
					resource.TestCheckResourceAttr("runscope_step.step", "variable.1.source", "response_status"),
					resource.TestCheckResourceAttr("runscope_step.step", "scripts.#", "2"),
					resource.TestCheckResourceAttr("runscope_step.step", "scripts.0", "log(\"script 1\");"),
					resource.TestCheckResourceAttr("runscope_step.step", "scripts.1", "log(\"script 2\");"),
					resource.TestCheckResourceAttr("runscope_step.step", "before_scripts.#", "1"),
					resource.TestCheckResourceAttr("runscope_step.step", "before_scripts.0", "log(\"before script\");"),
				),
			},
		},
	})
}

func TestAccStep_update(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckStepDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccStepDefaultConfig, bucketName, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStepExists("runscope_step.step"),
					resource.TestCheckResourceAttr("runscope_step.step", "step_type", "request"),
					resource.TestCheckResourceAttr("runscope_step.step", "skipped", "false"),
					resource.TestCheckResourceAttr("runscope_step.step", "note", ""),
					resource.TestCheckResourceAttr("runscope_step.step", "method", "GET"),
					resource.TestCheckResourceAttr("runscope_step.step", "body", ""),
					resource.TestCheckResourceAttr("runscope_step.step", "form_parameter.#", "0"),
					resource.TestCheckResourceAttr("runscope_step.step", "header.#", "0"),
					resource.TestCheckResourceAttr("runscope_step.step", "url", "https://example.org"),
					resource.TestCheckResourceAttr("runscope_step.step", "auth.#", "0"),
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.#", "0"),
					resource.TestCheckResourceAttr("runscope_step.step", "variable.#", "0"),
					resource.TestCheckResourceAttr("runscope_step.step", "scripts.#", "0"),
					resource.TestCheckResourceAttr("runscope_step.step", "before_scripts.#", "0"),
				),
			},
			{
				Config: fmt.Sprintf(testAccStepCustomConfig, bucketName, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStepExists("runscope_step.step"),
					resource.TestCheckResourceAttr("runscope_step.step", "step_type", "request"),
					resource.TestCheckResourceAttr("runscope_step.step", "skipped", "true"),
					resource.TestCheckResourceAttr("runscope_step.step", "note", "Step note"),
					resource.TestCheckResourceAttr("runscope_step.step", "method", "POST"),
					resource.TestCheckResourceAttr("runscope_step.step", "url", "https://example.org"),
					resource.TestCheckResourceAttr("runscope_step.step", "body", "Request body"),
					resource.TestCheckResourceAttr("runscope_step.step", "form_parameter.#", "3"),
					resource.TestCheckResourceAttr("runscope_step.step", "form_parameter.0.name", "a"),
					resource.TestCheckResourceAttr("runscope_step.step", "form_parameter.0.value", "1"),
					resource.TestCheckResourceAttr("runscope_step.step", "form_parameter.1.name", "a"),
					resource.TestCheckResourceAttr("runscope_step.step", "form_parameter.1.value", "2"),
					resource.TestCheckResourceAttr("runscope_step.step", "form_parameter.2.name", "b"),
					resource.TestCheckResourceAttr("runscope_step.step", "form_parameter.2.value", "3"),
					resource.TestCheckResourceAttr("runscope_step.step", "header.#", "3"),
					resource.TestCheckResourceAttr("runscope_step.step", "header.0.header", "Accept-Encoding"),
					resource.TestCheckResourceAttr("runscope_step.step", "header.0.value", "application/json"),
					resource.TestCheckResourceAttr("runscope_step.step", "header.1.header", "Accept-Encoding"),
					resource.TestCheckResourceAttr("runscope_step.step", "header.1.value", "application/xml"),
					resource.TestCheckResourceAttr("runscope_step.step", "header.2.header", "Authorization"),
					resource.TestCheckResourceAttr("runscope_step.step", "header.2.value", "Bearer bb74fe7b-b9f2-48bd-9445-bdc60e1edc6a"),
					resource.TestCheckResourceAttr("runscope_step.step", "auth.#", "1"),
					resource.TestCheckResourceAttr("runscope_step.step", "auth.0.username", "user"),
					resource.TestCheckResourceAttr("runscope_step.step", "auth.0.auth_type", "basic"),
					resource.TestCheckResourceAttr("runscope_step.step", "auth.0.password", "password1"),
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.#", "2"),
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.0.source", "response_status"),
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.0.comparison", "equal_number"),
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.0.value", "200"),
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.1.source", "response_json"),
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.1.comparison", "equal"),
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.1.value", "c5baeb4a-2379-478a-9cda-1b671de77cf9"),
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.1.property", "data.id"),
					resource.TestCheckResourceAttr("runscope_step.step", "variable.#", "2"),
					resource.TestCheckResourceAttr("runscope_step.step", "variable.0.name", "httpContentEncoding"),
					resource.TestCheckResourceAttr("runscope_step.step", "variable.0.source", "response_headers"),
					resource.TestCheckResourceAttr("runscope_step.step", "variable.0.property", "Content-Encoding"),
					resource.TestCheckResourceAttr("runscope_step.step", "variable.1.name", "httpStatus"),
					resource.TestCheckResourceAttr("runscope_step.step", "variable.1.source", "response_status"),
					resource.TestCheckResourceAttr("runscope_step.step", "scripts.#", "2"),
					resource.TestCheckResourceAttr("runscope_step.step", "scripts.0", "log(\"script 1\");"),
					resource.TestCheckResourceAttr("runscope_step.step", "scripts.1", "log(\"script 2\");"),
					resource.TestCheckResourceAttr("runscope_step.step", "before_scripts.#", "1"),
					resource.TestCheckResourceAttr("runscope_step.step", "before_scripts.0", "log(\"before script\");"),
				),
			},
		},
	})
}

func TestAccStep_multiple(t *testing.T) {
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

func TestAccStep_valid_variable_source(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckStepDestroy,
		Steps: func() []resource.TestStep {
			steps := make([]resource.TestStep, len(stepSources))
			for i, source := range stepSources {
				steps[i].Config = fmt.Sprintf(testAccStepVariableSourcesConfig, bucketName, teamId, source)
				steps[i].Check = resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("runscope_step.step", "variable.0.source", source),
				)
			}
			return steps
		}(),
	})
}

func TestAccStep_invalid_variable_source(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckStepDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccStepVariableSourcesConfig, bucketName, teamId, "invalid_source"),
				ExpectError: regexp.MustCompile("expected variable.0.source to be one of"),
			},
		},
	})
}

func TestAccStep_valid_assertion_source(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckStepDestroy,
		Steps: func() []resource.TestStep {
			steps := make([]resource.TestStep, len(stepSources))
			for i, source := range stepSources {
				steps[i].Config = fmt.Sprintf(testAccStepAssertionSourcesConfig, bucketName, teamId, source)
				steps[i].Check = resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.0.source", source),
				)
			}
			return steps
		}(),
	})
}

func TestAccStep_invalid_assertion_source(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckStepDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccStepAssertionSourcesConfig, bucketName, teamId, "invalid_source"),
				ExpectError: regexp.MustCompile("expected assertion.0.source to be one of"),
			},
		},
	})
}

func TestAccStep_valid_assertion_comparison(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckStepDestroy,
		Steps: func() []resource.TestStep {
			steps := make([]resource.TestStep, len(stepComparisons))
			for i, source := range stepComparisons {
				steps[i].Config = fmt.Sprintf(testAccStepAssertionComparisonsConfig, bucketName, teamId, source)
				steps[i].Check = resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("runscope_step.step", "assertion.0.comparison", source),
				)
			}
			return steps
		}(),
	})
}

func TestAccStep_invalid_assertion_comparison(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckStepDestroy,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testAccStepAssertionComparisonsConfig, bucketName, teamId, "invalid_compatison"),
				ExpectError: regexp.MustCompile("expected assertion.0.comparison to be one of"),
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

const testAccStepDefaultConfig = `
resource "runscope_bucket" "bucket" {
  name      = "%s"
  team_uuid = "%s"
}

resource "runscope_test" "test" {
  bucket_id   = runscope_bucket.bucket.id
  name        = "runscope test"
  description = "This is a test test..."
}

resource "runscope_step" "step" {
  bucket_id = runscope_bucket.bucket.id
  test_id   = runscope_test.test.id

  step_type = "request"
  method    = "GET"
  url       = "https://example.org"
}
`

const testAccStepVariableSourcesConfig = `
resource "runscope_bucket" "bucket" {
  name      = "%s"
  team_uuid = "%s"
}

resource "runscope_test" "test" {
  bucket_id   = runscope_bucket.bucket.id
  name        = "runscope test"
  description = "This is a test test..."
}

resource "runscope_step" "step" {
  bucket_id = runscope_bucket.bucket.id
  test_id   = runscope_test.test.id

  step_type = "request"
  method    = "GET"
  url       = "https://example.org"

  variable {
    name     = "httpStatus"
    source   = "%s"
  }
}
`

const testAccStepAssertionSourcesConfig = `
resource "runscope_bucket" "bucket" {
  name      = "%s"
  team_uuid = "%s"
}

resource "runscope_test" "test" {
  bucket_id   = runscope_bucket.bucket.id
  name        = "runscope test"
  description = "This is a test test..."
}

resource "runscope_step" "step" {
  bucket_id = runscope_bucket.bucket.id
  test_id   = runscope_test.test.id

  step_type = "request"
  method    = "GET"
  url       = "https://example.org"

  assertion {
    source     = "%s"
    comparison = "equal"
    value      = "c5baeb4a-2379-478a-9cda-1b671de77cf9"
    property   = "data.id"
  }
}
`

const testAccStepAssertionComparisonsConfig = `
resource "runscope_bucket" "bucket" {
  name      = "%s"
  team_uuid = "%s"
}

resource "runscope_test" "test" {
  bucket_id   = runscope_bucket.bucket.id
  name        = "runscope test"
  description = "This is a test test..."
}

resource "runscope_step" "step" {
  bucket_id = runscope_bucket.bucket.id
  test_id   = runscope_test.test.id

  step_type = "request"
  method    = "GET"
  url       = "https://example.org"

  assertion {
    source     = "response_status"
    comparison = "%s"
    property   = "data.id"
  }
}
`

const testAccStepCustomConfig = `
resource "runscope_bucket" "bucket" {
  name      = "%s"
  team_uuid = "%s"
}

resource "runscope_test" "test" {
  bucket_id   = runscope_bucket.bucket.id
  name        = "runscope test"
  description = "This is a test test..."
}

resource "runscope_step" "step" {
  bucket_id = runscope_bucket.bucket.id
  test_id   = runscope_test.test.id

  step_type = "request"
  method    = "POST"
  url       = "https://example.org"

  variable {
    name     = "httpStatus"
    source   = "response_status"
  }

  variable {
    name     = "httpContentEncoding"
    source   = "response_headers"
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
  
  header {
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

  skipped = true

  note = "Step note"

  body = "Request body"

  form_parameter {
    name = "a"
    value = "1"
  }

  form_parameter {
    name = "a"
    value = "2"
  }

  form_parameter {
    name = "b"
    value = "3"
  }
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
