package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRunscopeBucket(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceRunscopeBucketConfig, bucketName, teamId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.runscope_bucket.test", "name", bucketName),
				),
			},
		},
	})
}

const testAccDataSourceRunscopeBucketConfig = `
resource "runscope_bucket" "test" {
  name      = "%s"
  team_uuid = "%s"
}

data "runscope_bucket" "test" {
  key = runscope_bucket.test.id
}
`
