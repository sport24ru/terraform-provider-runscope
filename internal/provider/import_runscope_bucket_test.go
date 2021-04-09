package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportRunscopeBucket(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccImportRunscopeBucketConfig, bucketName, teamId),
			},
			{
				ResourceName:      "runscope_bucket.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccImportRunscopeBucketConfig = `
resource "runscope_bucket" "test" {
  name      = "%s"
  team_uuid = "%s"
}
`
