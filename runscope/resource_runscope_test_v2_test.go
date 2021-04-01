package runscope

import (
	"fmt"
	"os"
	"testing"

	"github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTestV2_basic(t *testing.T) {
	var test runscope.Test
	var defaultEnv runscope.Environment

	teamID := os.Getenv("RUNSCOPE_TEAM_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTestV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testRunscopeTestV2, teamID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTestV2Exists("runscope_test_v2.test", &test),
					testAccCheckEnvironmentExists("runscope_environment.environment", &defaultEnv),

					resource.TestCheckResourceAttr(
						"runscope_test_v2.test", "name", "runscope test"),
					resource.TestCheckResourceAttr(
						"runscope_test_v2.test", "description", "This is a test test..."),
					testAccCheckTestV2DefaultEnvId(&test, &defaultEnv),
				),
			},
		},
	})
}

func testAccCheckTestV2Destroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*runscope.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "runscope_test" {
			continue
		}

		_, err := client.ReadTest(&runscope.Test{ID: rs.Primary.ID, Bucket: &runscope.Bucket{Key: rs.Primary.Attributes["bucket_id"]}})

		if err == nil {
			return fmt.Errorf("Record %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTestV2Exists(n string, test *runscope.Test) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*runscope.Client)

		foundRecord, err := client.ReadTest(&runscope.Test{ID: rs.Primary.ID, Bucket: &runscope.Bucket{Key: rs.Primary.Attributes["bucket_id"]}})

		if err != nil {
			return err
		}

		if foundRecord.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*test = *foundRecord

		return nil
	}
}

func testAccCheckTestV2DefaultEnvId(test *runscope.Test, env *runscope.Environment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if test.DefaultEnvironmentID != env.ID {
			return fmt.Errorf("default environment ID is %s, %s expected", test.DefaultEnvironmentID, env.ID)
		}
		return nil
	}
}

const testRunscopeTestV2 = `
resource "runscope_bucket" "bucket" {
  name      = "terraform-provider-test"
  team_uuid = "%s"
}

resource "runscope_environment" "environment" {
  bucket_id = runscope_bucket.bucket.id
  name      = "test-default-environment"
}

resource "runscope_test_v2" "test" {
  bucket_id              = runscope_bucket.bucket.id
  name                   = "runscope test"
  description            = "This is a test test..."
  default_environment_id = runscope_environment.environment.id
}
`
