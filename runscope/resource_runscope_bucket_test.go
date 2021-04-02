package runscope

import (
	"context"
	"fmt"
	"github.com/terraform-providers/terraform-provider-runscope/internal/runscope"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBucket_basic(t *testing.T) {
	var bucket runscope.Bucket
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	bucketName := testAccRandomBucketName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testRunscopeBucketConfigA, bucketName, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists("runscope_bucket.bucket", &bucket),
					resource.TestCheckResourceAttr(
						"runscope_bucket.bucket", "name", bucketName),
				),
			},
		},
	})
}

func init() {
	resource.AddTestSweepers("runscope_bucket", &resource.Sweeper{
		Name: "runscope_bucket",
		F: func(region string) error {
			ctx := context.Background()

			client := runscope.NewClient(runscope.WithToken(os.Getenv("RUNSCOPE_ACCESS_TOKEN")))

			buckets, err := client.Bucket.List(ctx)
			if err != nil {
				return fmt.Errorf("Couldn't list bucket for sweeping")
			}

			for _, bucket := range buckets {
				if !(strings.HasPrefix(bucket.Name, testAccBucketNamePrefix) || bucket.Name == "terraform-provider-test") {
					continue
				}

				opts := &runscope.BucketDeleteOpts{}
				opts.Key = bucket.Key
				if err := client.Bucket.Delete(ctx, opts); err != nil {
					return err
				}
			}

			return nil
		},
	})
}

func testAccCheckBucketDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := testAccProvider.Meta().(*providerConfig).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "runscope_bucket" {
			continue
		}

		_, err := client.Bucket.Get(ctx, &runscope.BucketGetOpts{Key: rs.Primary.ID})

		if err == nil {
			return fmt.Errorf("Record %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckBucketExists(n string, bucket *runscope.Bucket) resource.TestCheckFunc {
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

		foundRecord, err := client.Bucket.Get(ctx, &runscope.BucketGetOpts{Key: rs.Primary.ID})

		if err != nil {
			return err
		}

		if foundRecord.Key != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		bucket = foundRecord

		return nil
	}
}

const testRunscopeBucketConfigA = `
resource "runscope_bucket" "bucket" {
  name      = "%s"
  team_uuid = "%s"
}`
