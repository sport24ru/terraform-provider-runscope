package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccEnvironment_create_default_shared_environment(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketId := testAccRandomBucketName()
	environment := runscope.Environment{}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			testAccEnvironmentDefaultConfigStep(testAccEnvironmentSharedDefaultConfig, bucketId, teamId, &environment),
		},
	})
}

func TestAccEnvironment_create_custom_shared_environment(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketId := testAccRandomBucketName()
	environment := runscope.Environment{}

	recipientId, recipientIdOk := os.LookupEnv("RUNSCOPE_RECIPIENT_ID")
	recipientName, recipientNameOk := os.LookupEnv("RUNSCOPE_RECIPIENT_NAME")
	recipientEmail, recipientEmailOk := os.LookupEnv("RUNSCOPE_RECIPIENT_EMAIL")

	if !(recipientIdOk && recipientNameOk && recipientEmailOk) {
		t.Skip("All of RUNSCOPE_RECIPIENT_ID, RUNSCOPE_RECIPIENT_NAME, RUNSCOPE_RECIPIENT_EMAIL should be set")
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			testAccEnvironmentCustomConfigStep(testAccEnvironmentSharedCustomConfig, bucketId, teamId, recipientId, recipientName, recipientEmail, &environment),
		},
	})
}

func TestAccEnvironment_update_custom_shared_environment(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketId := testAccRandomBucketName()
	environment := runscope.Environment{}

	recipientId, recipientIdOk := os.LookupEnv("RUNSCOPE_RECIPIENT_ID")
	recipientName, recipientNameOk := os.LookupEnv("RUNSCOPE_RECIPIENT_NAME")
	recipientEmail, recipientEmailOk := os.LookupEnv("RUNSCOPE_RECIPIENT_EMAIL")

	if !(recipientIdOk && recipientNameOk && recipientEmailOk) {
		t.Skip("All of RUNSCOPE_RECIPIENT_ID, RUNSCOPE_RECIPIENT_NAME, RUNSCOPE_RECIPIENT_EMAIL should be set")
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			testAccEnvironmentDefaultConfigStep(testAccEnvironmentSharedDefaultConfig, bucketId, teamId, &environment),
			testAccEnvironmentCustomConfigStep(testAccEnvironmentSharedCustomConfig, bucketId, teamId, recipientId, recipientName, recipientEmail, &environment),
		},
	})
}

func TestAccEnvironment_create_default_test_environment(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketId := testAccRandomBucketName()
	environment := runscope.Environment{}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			testAccEnvironmentDefaultConfigStep(testAccEnvironmentTestDefaultConfig, bucketId, teamId, &environment),
		},
	})
}

func TestAccEnvironment_create_custom_test_environment(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketId := testAccRandomBucketName()
	environment := runscope.Environment{}

	recipientId, recipientIdOk := os.LookupEnv("RUNSCOPE_RECIPIENT_ID")
	recipientName, recipientNameOk := os.LookupEnv("RUNSCOPE_RECIPIENT_NAME")
	recipientEmail, recipientEmailOk := os.LookupEnv("RUNSCOPE_RECIPIENT_EMAIL")

	if !(recipientIdOk && recipientNameOk && recipientEmailOk) {
		t.Skip("All of RUNSCOPE_RECIPIENT_ID, RUNSCOPE_RECIPIENT_NAME, RUNSCOPE_RECIPIENT_EMAIL should be set")
		return
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			testAccEnvironmentCustomConfigStep(testAccEnvironmentTestCustomConfig, bucketId, teamId, recipientId, recipientName, recipientEmail, &environment),
		},
	})
}

func TestAccEnvironment_update_custom_test_environment(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketId := testAccRandomBucketName()
	environment := runscope.Environment{}

	recipientId, recipientIdOk := os.LookupEnv("RUNSCOPE_RECIPIENT_ID")
	recipientName, recipientNameOk := os.LookupEnv("RUNSCOPE_RECIPIENT_NAME")
	recipientEmail, recipientEmailOk := os.LookupEnv("RUNSCOPE_RECIPIENT_EMAIL")

	if !(recipientIdOk && recipientNameOk && recipientEmailOk) {
		t.Skip("All of RUNSCOPE_RECIPIENT_ID, RUNSCOPE_RECIPIENT_NAME, RUNSCOPE_RECIPIENT_EMAIL should be set")
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			testAccEnvironmentDefaultConfigStep(testAccEnvironmentTestDefaultConfig, bucketId, teamId, &environment),
			testAccEnvironmentCustomConfigStep(testAccEnvironmentTestCustomConfig, bucketId, teamId, recipientId, recipientName, recipientEmail, &environment),
		},
	})
}

func TestAccEnvironment_create_nested_test_environment(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketId := testAccRandomBucketName()
	environment := runscope.Environment{}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			testAccEnvironmentDefaultConfigStep(testAccEnvironmentTestNestedConfig, bucketId, teamId, &environment),
		},
	})
}

func testAccCheckEnvironmentDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := testAccProvider.Meta().(*providerConfig).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "runscope_environment" {
			continue
		}

		opts := runscope.EnvironmentDeleteOpts{}
		opts.Id = rs.Primary.ID
		opts.BucketId = rs.Primary.Attributes["bucket_id"]
		opts.TestId = rs.Primary.Attributes["test_id"]

		if err := client.Environment.Delete(ctx, &opts); err == nil {
			return fmt.Errorf("record %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckEnvironmentExists(n string, e *runscope.Environment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		client := testAccProvider.Meta().(*providerConfig).client

		opts := runscope.EnvironmentGetOpts{}
		opts.Id = rs.Primary.ID
		opts.BucketId = rs.Primary.Attributes["bucket_id"]
		opts.TestId = rs.Primary.Attributes["test_id"]

		environment, err := client.Environment.Get(ctx, &opts)
		if err != nil {
			return err
		}

		if environment.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		*e = *environment

		return nil
	}
}

func testAccCheckEnvironmentRecipient(e *runscope.Environment, expectedId string, expectedName string, expectedEmail string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		id := e.Emails.Recipients[0].Id
		if id != expectedId {
			return fmt.Errorf("expected recipient ID '%s', got '%s'", expectedId, id)
		}
		name := e.Emails.Recipients[0].Name
		if name != expectedName {
			return fmt.Errorf("expected recipient name '%s', got '%s'", expectedName, name)
		}
		email := e.Emails.Recipients[0].Email
		if email != expectedEmail {
			return fmt.Errorf("expected recipient email '%s', got '%s'", expectedEmail, email)
		}

		return nil
	}
}

func testAccEnvironmentDefaultConfigStep(config, bucketId, teamId string, environment *runscope.Environment) resource.TestStep {
	return resource.TestStep{
		Config: fmt.Sprintf(config, bucketId, teamId),
		Check: resource.ComposeTestCheckFunc(
			testAccCheckEnvironmentExists("runscope_environment.environment", environment),
			resource.TestCheckResourceAttr("runscope_environment.environment", "name", "environment"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "script", ""),
			resource.TestCheckResourceAttr("runscope_environment.environment", "preserve_cookies", "false"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "initial_variables.%", "0"),
			resource.TestCheckNoResourceAttr("runscope_environment.environment", "integrations"),
			resource.TestCheckNoResourceAttr("runscope_environment.environment", "regions"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "retry_on_failure", "false"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "stop_on_failure", "false"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "verify_ssl", "true"),
			resource.TestCheckNoResourceAttr("runscope_environment.environment", "webhooks"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "email.#", "0"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "parent_environment_id", ""),
			resource.TestCheckResourceAttr("runscope_environment.environment", "client_certificate", ""),
		),
	}
}

func testAccEnvironmentCustomConfigStep(config, bucketId, teamId, recipientId, recipientName, recipientEmail string, environment *runscope.Environment) resource.TestStep {
	return resource.TestStep{
		Config: fmt.Sprintf(config, bucketId, teamId, recipientId),
		Check: resource.ComposeTestCheckFunc(
			testAccCheckEnvironmentExists("runscope_environment.environment", environment),
			resource.TestCheckResourceAttr("runscope_environment.environment", "name", "environment"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "script", "1;"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "preserve_cookies", "true"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "initial_variables.%", "2"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "integrations.#", "1"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "regions.#", "2"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "retry_on_failure", "true"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "stop_on_failure", "true"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "verify_ssl", "true"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "webhooks.#", "1"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "email.#", "1"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "email.0.notify_all", "true"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "email.0.notify_on", "all"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "email.0.notify_threshold", "1"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "email.0.recipient.#", "1"),
			resource.TestCheckResourceAttr("runscope_environment.environment", "email.0.recipient.0.id", recipientId),
			testAccCheckEnvironmentRecipient(environment, recipientId, recipientName, recipientEmail),
			resource.TestCheckResourceAttr("runscope_environment.environment", "parent_environment_id", ""),
			resource.TestCheckResourceAttr("runscope_environment.environment", "client_certificate", testAccEnvironmentClientCertficate),
		),
	}
}

const testAccEnvironmentClientCertficate = `-----BEGIN CERTIFICATE-----
MIIDDTCCAfWgAwIBAgIUd8JoBoWhPUHSqxgMvDqTgBFmHTswDQYJKoZIhvcNAQEL
BQAwFjEUMBIGA1UECwwLZXhhbXBsZS5vcmcwHhcNMjEwNDExMjEyMzE5WhcNMzEw
NDA5MjEyMzE5WjAWMRQwEgYDVQQLDAtleGFtcGxlLm9yZzCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBAK1FEXUn5I85YIxc88ql7T4Vn1nIdcKsrGfMR3oH
ZEFTXP/TG0GCgAPLgCmBCLJZUJVcgpENiIzqO+JoZ0daBm5Cf6Y/ZZFX/VXxZtSV
hsLnNozf3IXl5T00JXPg2JYTSqZZfBbAREQQAZuucsSgP4t7kP0Q9L/fiCUkEGRe
jU4oncwRnyNv85qf6rRHtK0+REdvMm56oVHqvXR4k5EvtEpf1qOfeeuJ+ZCh/0yu
zUarhY5jwdXircRiDmWfkk3PhqP5lsBiaDbXemLb0DDWDBzFGj9aebptp0bfSZ/0
KsQdYD4MM0eDep9a4JvT0nxXZ+RvWb3o081i+pz/AcT85nkCAwEAAaNTMFEwHQYD
VR0OBBYEFLzbHqu0i6Kkz2kto5olS3nmR9OEMB8GA1UdIwQYMBaAFLzbHqu0i6Kk
z2kto5olS3nmR9OEMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEB
ADmmcH+Ie8KscvEuNm5kWVCidixtm31atlFtLUHLRn9898e5k5FviGY13RS9NAJB
0h5PLwIytm2+FkjBZ1hukWCQgijrKKKfdcKzes6MQglGru8RqUqlP9r/T2ly+YpM
D6vzUeIN4QT1fxlR017n8OFATsifdaFMBezST84wIoGWSf+njIqsgGpMhn/8/Xb3
5cQ9MngrRvpmwlBNgujPNo08X9MxduLX4Yz19jeXufG72ebrnriTxRxJNTW+j/Pt
JlalxgJ13chSeQkb2U9/1Es4PGfYvNXIJ6YmOeT8O6CIt+cHK5yAiarKJE59GJPJ
upwBSEp4QTHSUGFgnR6PZfQ=
-----END CERTIFICATE-----
`

const testAccEnvironmentSharedDefaultConfig = `
resource "runscope_bucket" "bucket" {
	name      = "%s"
	team_uuid = "%s"
}

resource "runscope_environment" "environment" {
	bucket_id = runscope_bucket.bucket.id
	name      = "environment"
}
`

const testAccEnvironmentSharedCustomConfig = `
resource "runscope_bucket" "bucket" {
	name      = "%[1]s"
	team_uuid = "%[2]s"
}

data "runscope_integration" "slack" {
	team_uuid = "%[2]s"
	type      = "slack"
}

resource "runscope_environment" "environment" {
  bucket_id = runscope_bucket.bucket.id
  name      = "environment"
  
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
  stop_on_failure  = true
  webhooks         = ["https://example.com"]
  
  email {
  	notify_all       = true
  	notify_on        = "all"
  	notify_threshold = 1
  	recipient {
      id = "%[3]s"
    }
  }

  client_certificate = <<EOF
` + testAccEnvironmentClientCertficate + `EOF
}
`

const testAccEnvironmentTestDefaultConfig = `
resource "runscope_bucket" "bucket" {
  name      = "%s"
  team_uuid = "%s"
}

resource "runscope_test" "test" {
  bucket_id   = runscope_bucket.bucket.id
  name        = "runscope test"
  description = "description"
}

resource "runscope_environment" "environment" {
  bucket_id = runscope_bucket.bucket.id
  test_id   = runscope_test.test.id
  name      = "environment"
}
`

const testAccEnvironmentTestCustomConfig = `
resource "runscope_bucket" "bucket" {
  name      = "%[1]s"
  team_uuid = "%[2]s"
}

resource "runscope_test" "test" {
  bucket_id   = runscope_bucket.bucket.id
  name        = "runscope test"
  description = "description"
}

data "runscope_integration" "slack" {
  team_uuid = "%[2]s"
  type      = "slack"
}

resource "runscope_environment" "environment" {
  bucket_id = runscope_bucket.bucket.id
  test_id   = runscope_test.test.id
  name      = "environment"
      
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
  stop_on_failure  = true
  webhooks         = ["https://example.com"]
      
  email {
  	notify_all       = true
  	notify_on        = "all"
  	notify_threshold = 1
  	recipient {
      id = "%[3]s"
    }
  }

  client_certificate = <<EOF
` + testAccEnvironmentClientCertficate + `EOF
}
`

const testAccEnvironmentTestNestedConfig = `
resource "runscope_bucket" "bucket" {
  name      = "%s"
  team_uuid = "%s"
}

resource "runscope_test" "test" {
  bucket_id   = runscope_bucket.bucket.id
  name        = "runscope test"
  description = "description"
}

resource "runscope_environment" "environment" {
  bucket_id = runscope_bucket.bucket.id
  test_id   = runscope_test.test.id
  name      = "environment"
  regions   = ["us1", "eu1"]
}

resource "runscope_environment" "environment_child" {
  bucket_id             = runscope_bucket.bucket.id
  test_id               = runscope_test.test.id
  parent_environment_id = runscope_environment.environment.id
  name                  = "environment-child"
}
`
