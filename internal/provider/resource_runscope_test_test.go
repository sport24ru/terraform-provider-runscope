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

func TestAccTest_create_default_test(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()
	test := &runscope.Test{}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTestDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTestDefaultConfig, bucketName, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTestExists("runscope_test.test", test),
					resource.TestCheckResourceAttr("runscope_test.test", "name", "runscope test"),
					resource.TestCheckResourceAttr("runscope_test.test", "description", ""),
					resource.TestCheckResourceAttrSet("runscope_test.test", "default_environment_id"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "created_at"),
					resource.TestCheckResourceAttr("runscope_test.test", "created_by.#", "1"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "created_by.0.id"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "created_by.0.name"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "created_by.0.email"),
				),
			},
		},
	})
}

func TestAccTest_create_custom_test(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()
	test := &runscope.Test{}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTestDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTestCustomConfig, bucketName, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTestExists("runscope_test.test", test),
					resource.TestCheckResourceAttr("runscope_test.test", "name", "runscope custom test"),
					resource.TestCheckResourceAttr("runscope_test.test", "description", "runscope custom test description"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "default_environment_id"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "created_at"),
					resource.TestCheckResourceAttr("runscope_test.test", "created_by.#", "1"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "created_by.0.id"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "created_by.0.name"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "created_by.0.email"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "trigger_url"),
				),
			},
			{
				ResourceName:      "runscope_test.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["runscope_test.test"]
					if !ok {
						return "", fmt.Errorf("not found runscope_test.test")
					}
					return fmt.Sprintf("%s/%s", rs.Primary.Attributes["bucket_id"], rs.Primary.ID), nil
				},
			},
		},
	})
}

func TestAccTest_update_test(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()
	test1 := &runscope.Test{}
	test2 := &runscope.Test{}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTestDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTestDefaultConfig, bucketName, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTestExists("runscope_test.test", test1),
					resource.TestCheckResourceAttr("runscope_test.test", "name", "runscope test"),
					resource.TestCheckResourceAttr("runscope_test.test", "description", ""),
					resource.TestCheckResourceAttrSet("runscope_test.test", "default_environment_id"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "created_at"),
					resource.TestCheckResourceAttr("runscope_test.test", "created_by.#", "1"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "created_by.0.id"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "created_by.0.name"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "created_by.0.email"),
				),
			},
			{
				Config: fmt.Sprintf(testAccTestCustomConfig, bucketName, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTestExists("runscope_test.test", test2),
					resource.TestCheckResourceAttr("runscope_test.test", "name", "runscope custom test"),
					resource.TestCheckResourceAttr("runscope_test.test", "description", "runscope custom test description"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "default_environment_id"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "created_at"),
					resource.TestCheckResourceAttr("runscope_test.test", "created_by.#", "1"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "created_by.0.id"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "created_by.0.name"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "created_by.0.email"),
					resource.TestCheckResourceAttrSet("runscope_test.test", "trigger_url"),
					testAccCheckTestIdEqual(test1, test2),
				),
			},
		},
	})
}

func testAccCheckTestDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := testAccProvider.Meta().(*providerConfig).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "runscope_test" {
			continue
		}

		opts := runscope.TestDeleteOpts{}
		opts.Id = rs.Primary.ID
		opts.BucketId = rs.Primary.Attributes["bucket_id"]

		if err := client.Test.Delete(ctx, opts); err == nil {
			return fmt.Errorf("record %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTestExists(n string, t *runscope.Test) resource.TestCheckFunc {
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

		opts := runscope.TestGetOpts{}
		opts.Id = rs.Primary.ID
		opts.BucketId = rs.Primary.Attributes["bucket_id"]
		test, err := client.Test.Get(ctx, opts)
		if err != nil {
			return err
		}

		if test.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		*t = *test

		return nil
	}
}

func testAccCheckTestIdEqual(t1, t2 *runscope.Test) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if t1.Id != t2.Id {
			return fmt.Errorf("expected than \"%s\" equal to \"%s\"", t1.Id, t2.Id)
		}
		return nil
	}
}

const testAccTestDefaultConfig = `
resource "runscope_bucket" "bucket" {
  name      = "%s"
  team_uuid = "%s"
}

resource "runscope_test" "test" {
  bucket_id = runscope_bucket.bucket.id

  name = "runscope test"
}
`

const testAccTestCustomConfig = `
resource "runscope_bucket" "bucket" {
  name      = "%s"
  team_uuid = "%s"
}

resource "runscope_test" "test" {
  bucket_id = runscope_bucket.bucket.id

  name        = "runscope custom test"
  description = "runscope custom test description"
}
`
