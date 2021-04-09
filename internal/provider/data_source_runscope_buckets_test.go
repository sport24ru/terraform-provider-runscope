package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRunscopeBuckets(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceRunscopeBucketsConfig, teamId, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.runscope_buckets.test", "keys.#", "1"),
				),
			},
		},
	})
}

const testAccDataSourceRunscopeBucketsConfig = `
resource "runscope_bucket" "test" {
  team_uuid = "%s"
  name      = "%s"
}

data "runscope_buckets" "test" {
  filter {
    name   = "name"
    values = [runscope_bucket.test.name]
  }
}
`
