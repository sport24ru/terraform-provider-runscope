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

func TestAccSchedule_create_default_schedule(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()
	schedule := &runscope.Schedule{}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccScheduleDefaultConfig, bucketName, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScheduleExists("runscope_schedule.daily", schedule),
					resource.TestCheckResourceAttr("runscope_schedule.daily", "interval", "1d"),
					resource.TestCheckResourceAttr("runscope_schedule.daily", "note", ""),
					resource.TestCheckResourceAttrSet("runscope_schedule.daily", "exported_at"),
				),
			},
		},
	})
}

func TestAccSchedule_create_custom_schedule(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()
	schedule := &runscope.Schedule{}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccScheduleCustomConfig, bucketName, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScheduleExists("runscope_schedule.daily", schedule),
					resource.TestCheckResourceAttr("runscope_schedule.daily", "interval", "6h"),
					resource.TestCheckResourceAttr("runscope_schedule.daily", "note", "schedule note"),
					resource.TestCheckResourceAttrSet("runscope_schedule.daily", "exported_at"),
				),
			},
		},
	})
}

func TestAccSchedule_update_schedule(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()
	schedule := &runscope.Schedule{}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccScheduleDefaultConfig, bucketName, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScheduleExists("runscope_schedule.daily", schedule),
					resource.TestCheckResourceAttr("runscope_schedule.daily", "interval", "1d"),
					resource.TestCheckResourceAttr("runscope_schedule.daily", "note", ""),
					resource.TestCheckResourceAttrSet("runscope_schedule.daily", "exported_at"),
				),
			},
			{
				Config: fmt.Sprintf(testAccScheduleCustomConfig, bucketName, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScheduleExists("runscope_schedule.daily", schedule),
					resource.TestCheckResourceAttr("runscope_schedule.daily", "interval", "6h"),
					resource.TestCheckResourceAttr("runscope_schedule.daily", "note", "schedule note"),
					resource.TestCheckResourceAttrSet("runscope_schedule.daily", "exported_at"),
				),
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

func testAccCheckScheduleExists(n string, sch *runscope.Schedule) resource.TestCheckFunc {
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
		if schedule.EnvironmentId != rs.Primary.Attributes["environment_id"] {
			return fmt.Errorf("Expected environment Id `%s`, got `%s`", schedule.EnvironmentId, rs.Primary.Attributes["environment_id"])
		}

		*sch = *schedule

		return nil
	}
}

const testAccScheduleDefaultConfig = `
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
}

resource "runscope_environment" "environment2" {
  bucket_id = runscope_bucket.bucket.id
  name      = "test-environment2"
}

resource "runscope_schedule" "daily" {
  bucket_id      = runscope_bucket.bucket.id
  test_id        = runscope_test.test.id
  environment_id = runscope_environment.environment.id
  interval       = "1d"
}
`

const testAccScheduleCustomConfig = `
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
}

resource "runscope_environment" "environment2" {
  bucket_id = runscope_bucket.bucket.id
  name      = "test-environment2"
}

resource "runscope_schedule" "daily" {
  bucket_id      = runscope_bucket.bucket.id
  test_id        = runscope_test.test.id
  environment_id = runscope_environment.environment2.id
  interval       = "6h"
  note           = "schedule note"
}
`
