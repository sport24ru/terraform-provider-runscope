package runscope

import (
	"fmt"
	"os"
	"testing"

	"github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccEnvironment_basic(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	environment := runscope.Environment{}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testRunscopeEnvironmentMinimal, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists("runscope_environment.environmentA", &environment),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "name", "test-environment"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "script", ""),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "preserve_cookies", "false"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "initial_variables.%", "0"),
					resource.TestCheckNoResourceAttr("runscope_environment.environmentA", "integrations"),
					resource.TestCheckNoResourceAttr("runscope_environment.environmentA", "regions"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "retry_on_failure", "false"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "verify_ssl", "true"),
					resource.TestCheckNoResourceAttr("runscope_environment.environmentA", "webhooks"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "email.#", "0"),
				),
			},
			{
				Config: fmt.Sprintf(testRunscopeEnvironmentFull, teamId, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists("runscope_environment.environmentA", &environment),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "name", "test-environment"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "script", "1;"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "preserve_cookies", "true"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "initial_variables.%", "2"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "integrations.#", "1"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "regions.#", "2"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "retry_on_failure", "true"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "verify_ssl", "true"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "webhooks.#", "1"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "email.#", "0"),
				),
			},
		},
	})
}

func TestAccEnvironment_email(t *testing.T) {
	teamID := os.Getenv("RUNSCOPE_TEAM_ID")
	environment := runscope.Environment{}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testRunscopeEnvironmentConfigWithEmail, teamID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists("runscope_environment.environmentA", &environment),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "name", "test-environment"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "verify_ssl", "true"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "email.#", "1"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "email.0.notify_all", "true"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "email.0.notify_on", "all"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "email.0.notify_threshold", "1"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "email.0.recipient.#", "0"),
				),
			},
		},
	})
}

func TestAccEnvironment_update_email(t *testing.T) {
	teamID := os.Getenv("RUNSCOPE_TEAM_ID")
	recipientId, recipientIdOk := os.LookupEnv("RUNSCOPE_RECIPIENT_ID")
	if !recipientIdOk {
		t.Skip("RUNSCOPE_RECIPIENT_ID should be set")
		return
	}

	environment := runscope.Environment{}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testRunscopeEnvironmentMinimal, teamID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists("runscope_environment.environmentA", &environment),
					testAccCheckEnvironmentEmail(&environment, false, "", 0, 0),
				),
			},
			{
				Config: fmt.Sprintf(testRunscopeEnvironmentConfigWithEmailRecipient, teamID, recipientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists("runscope_environment.environmentA", &environment),
					testAccCheckEnvironmentEmail(&environment, true, "all", 1, 1),
				),
			},
			{
				Config: fmt.Sprintf(testRunscopeEnvironmentMinimal, teamID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists("runscope_environment.environmentA", &environment),
					testAccCheckEnvironmentEmail(&environment, false, "all", 1, 0),
				),
			},
		},
	})
}

func TestAccEnvironment_email_recipient(t *testing.T) {
	teamId, ok := os.LookupEnv("RUNSCOPE_TEAM_ID")
	if !ok {
		t.Skip("RUNSCOPE_TEAM_ID should be set")
		return
	}

	recipientId, recipientIdOk := os.LookupEnv("RUNSCOPE_RECIPIENT_ID")
	recipientName, recipientNameOk := os.LookupEnv("RUNSCOPE_RECIPIENT_NAME")
	recipientEmail, recipientEmailOk := os.LookupEnv("RUNSCOPE_RECIPIENT_EMAIL")

	if !(recipientIdOk && recipientNameOk && recipientEmailOk) {
		t.Skip("All of RUNSCOPE_RECIPIENT_ID, RUNSCOPE_RECIPIENT_NAME, RUNSCOPE_RECIPIENT_EMAIL should be set")
		return
	}

	hash := recipientsHash(map[string]interface{}{
		"id": recipientId,
	})

	environment := runscope.Environment{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testRunscopeEnvironmentConfigWithEmailRecipient, teamId, recipientId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists("runscope_environment.environmentA", &environment),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "name", "test-environment"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "email.#", "1"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "email.0.notify_all", "true"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "email.0.notify_on", "all"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "email.0.notify_threshold", "1"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", "email.0.recipient.#", "1"),
					resource.TestCheckResourceAttr("runscope_environment.environmentA", fmt.Sprintf("email.0.recipient.%d.id", hash), recipientId),
					testAccCheckEnvironmentRecipient(&environment, recipientId, recipientName, recipientEmail),
				),
			},
		},
	})
}

func TestAccEnvironment_do_not_verify_ssl(t *testing.T) {
	teamID := os.Getenv("RUNSCOPE_TEAM_ID")
	environment := runscope.Environment{}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testRunscopeEnvironmentConfigB, teamID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists("runscope_environment.environmentB", &environment),
					resource.TestCheckResourceAttr(
						"runscope_environment.environmentB", "name", "test-no-ssl"),
					resource.TestCheckResourceAttr(
						"runscope_environment.environmentB", "verify_ssl", "false")),
			},
		},
	})
}

func testAccCheckEnvironmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*runscope.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "runscope_environment" {
			continue
		}

		var err error
		bucketID := rs.Primary.Attributes["bucket_id"]
		testID := rs.Primary.Attributes["test_id"]
		if testID != "" {
			err = client.DeleteEnvironment(&runscope.Environment{ID: rs.Primary.ID},
				&runscope.Bucket{Key: bucketID})
		} else {
			err = client.DeleteEnvironment(&runscope.Environment{ID: rs.Primary.ID},
				&runscope.Bucket{Key: bucketID})
		}

		if err == nil {
			return fmt.Errorf("Record %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckEnvironmentExists(n string, e *runscope.Environment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*runscope.Client)

		var foundRecord *runscope.Environment
		var err error

		environment := new(runscope.Environment)
		environment.ID = rs.Primary.ID
		bucketID := rs.Primary.Attributes["bucket_id"]
		testID := rs.Primary.Attributes["test_id"]
		if testID != "" {
			foundRecord, err = client.ReadTestEnvironment(environment,
				&runscope.Test{
					ID:     testID,
					Bucket: &runscope.Bucket{Key: bucketID}})
		} else {
			foundRecord, err = client.ReadSharedEnvironment(environment,
				&runscope.Bucket{Key: bucketID})
		}

		if err != nil {
			return err
		}

		if foundRecord.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*e = *foundRecord

		return nil
	}
}

func testAccCheckEnvironmentEmail(e *runscope.Environment, expectedNotifyAll bool, expectedNotifyOn string, expectedNotifyThreshold int, expectedNumRecipients int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if e.EmailSettings.NotifyAll != expectedNotifyAll {
			return fmt.Errorf("Expected NotifyAll '%v', got '%v'", expectedNotifyAll, e.EmailSettings.NotifyAll)
		}
		if e.EmailSettings.NotifyOn != expectedNotifyOn {
			return fmt.Errorf("Expected NotifyOn '%s', got '%s'", expectedNotifyOn, e.EmailSettings.NotifyOn)
		}
		if e.EmailSettings.NotifyThreshold != expectedNotifyThreshold {
			return fmt.Errorf("Expected NotifyThreshold '%d', got '%d'", expectedNotifyThreshold, e.EmailSettings.NotifyThreshold)
		}
		if len(e.EmailSettings.Recipients) != expectedNumRecipients {
			return fmt.Errorf("Expected '%d' recipients, got '%d'", expectedNumRecipients, len(e.EmailSettings.Recipients))
		}
		return nil
	}
}

func testAccCheckEnvironmentRecipient(e *runscope.Environment, expectedId string, expectedName string, expectedEmail string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		id := e.EmailSettings.Recipients[0].ID
		if id != expectedId {
			return fmt.Errorf("Expected recipient ID '%s', got '%s'", expectedId, id)
		}
		name := e.EmailSettings.Recipients[0].Name
		if name != expectedName {
			return fmt.Errorf("Expected recipient name '%s', got '%s'", expectedName, name)
		}
		email := e.EmailSettings.Recipients[0].Email
		if email != expectedEmail {
			return fmt.Errorf("Expected recipient email '%s', got '%s'", expectedEmail, email)
		}

		return nil
	}
}

const testRunscopeEnvironmentMinimal = `
resource "runscope_bucket" "bucket" {
	name      = "terraform-provider-test"
	team_uuid = "%s"
}

resource "runscope_environment" "environmentA" {
	bucket_id = runscope_bucket.bucket.id
	name      = "test-environment"
}
`

const testRunscopeEnvironmentFull = `
resource "runscope_test" "test" {
	bucket_id = runscope_bucket.bucket.id
	name = "runscope test"
	description = "This is a test test..."
}

resource "runscope_bucket" "bucket" {
	name      = "terraform-provider-test"
	team_uuid = "%s"
}

data "runscope_integration" "slack" {
	team_uuid = "%s"
	type      = "slack"
}

resource "runscope_environment" "environmentA" {
	bucket_id    = runscope_bucket.bucket.id
	name         = "test-environment"

	script = "1;"

	preserve_cookies = true

	initial_variables = {
		var1 = "true"
		var2 = "value2"
	}

	integrations = [
		data.runscope_integration.slack.id,
	]

	regions = ["us1", "eu1"]

	remote_agent {
		name = "test agent"
		uuid = "arbitrary-string"
	}

	retry_on_failure = true
	webhooks         = ["https://example.com"]
}
`

const testRunscopeEnvironmentConfigWithEmail = `
resource "runscope_bucket" "bucket" {
	name      = "terraform-provider-test"
	team_uuid = "%s"
}

resource "runscope_environment" "environmentA" {
	bucket_id    = runscope_bucket.bucket.id
	name         = "test-environment"

	email {
		notify_all       = true
		notify_on        = "all"
		notify_threshold = 1
	}
}
`

const testRunscopeEnvironmentConfigWithEmailRecipient = `
resource "runscope_bucket" "bucket" {
  name      = "terraform-provider-test"
  team_uuid = "%s"
}

resource "runscope_environment" "environmentA" {
	bucket_id    = "${runscope_bucket.bucket.id}"
	name         = "test-environment"

	email {
		notify_all       = true
		notify_on        = "all"
		notify_threshold = 1
		recipient {
 			id = "%s"
        }
	}
}
`

const testRunscopeEnvironmentConfigB = `
resource "runscope_bucket" "bucket" {
	name      = "terraform-provider-test"
	team_uuid = "%s"
}

resource "runscope_environment" "environmentB" {
	bucket_id = "${runscope_bucket.bucket.id}"
	name      = "test-no-ssl"

	verify_ssl = false
}
`
