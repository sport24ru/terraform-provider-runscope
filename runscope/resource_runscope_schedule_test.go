package runscope

import (
	"context"
	"fmt"
	"github.com/terraform-providers/terraform-provider-runscope/internal/runscope"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSchedule_basic(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testRunscopeSchedule, bucketName, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScheduleExists("runscope_schedule.daily"),
					resource.TestCheckResourceAttr(
						"runscope_schedule.daily", "note", "This is a daily schedule"),
					resource.TestCheckResourceAttr(
						"runscope_schedule.daily", "interval", "1d")),
			},
		},
	})
}

func testAccCheckScheduleDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := testAccProvider.Meta().(*providerConfig).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "runscope_schedule" {
			continue
		}

		opts := &runscope.ScheduleDeleteOpts{}
		opts.Id = rs.Primary.ID
		opts.BucketId = rs.Primary.Attributes["bucket_id"]
		opts.TestId = rs.Primary.Attributes["test_id"]

		if err := client.Schedule.Delete(ctx, opts); err == nil {
			return fmt.Errorf("Record %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckScheduleExists(n string) resource.TestCheckFunc {
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

		opts := &runscope.ScheduleGetOpts{}
		opts.Id = rs.Primary.ID
		opts.BucketId = rs.Primary.Attributes["bucket_id"]
		opts.TestId = rs.Primary.Attributes["test_id"]

		schedule, err := client.Schedule.Get(ctx, opts)
		if err != nil {
			return err
		}

		if schedule.Id != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		return nil
	}
}

const testRunscopeSchedule = `
resource "runscope_bucket" "bucket" {
  name      = "%s"
  team_uuid = "%s"
}

resource "runscope_test" "test" {
  bucket_id   = runscope_bucket.bucket.id
  name        = "runscope test"
  description = "This is a test test..."
}

resource "runscope_environment" "environment" {
  bucket_id = runscope_bucket.bucket.id
  name      = "test-environment"

  initial_variables = {
    var1 = "true"
    var2 = "value2"
  }
}

resource "runscope_schedule" "daily" {
  bucket_id      = runscope_bucket.bucket.id
  test_id        = runscope_test.test.id
  environment_id = runscope_environment.environment.id
  interval       = "1d"
  note           = "This is a daily schedule"
}
`
