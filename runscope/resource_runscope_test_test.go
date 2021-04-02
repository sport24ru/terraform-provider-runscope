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

func TestAccTest_basic(t *testing.T) {
	var test runscope.Test
	teamID := os.Getenv("RUNSCOPE_TEAM_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTestDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testRunscopeTestConfigA, teamID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTestExists("runscope_test.test", &test),
					resource.TestCheckResourceAttr(
						"runscope_test.test", "name", "runscope test"),
					resource.TestCheckResourceAttr(
						"runscope_test.test", "description", "This is a test test..."),
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

const testRunscopeTestConfigA = `
resource "runscope_test" "test" {
  bucket_id = "${runscope_bucket.bucket.id}"
  name = "runscope test"
  description = "This is a test test..."
}

resource "runscope_bucket" "bucket" {
  name = "terraform-provider-test"
  team_uuid = "%s"
}
`
