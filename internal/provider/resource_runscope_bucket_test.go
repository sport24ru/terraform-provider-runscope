package provider

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("runscope_bucket", &resource.Sweeper{
		Name: "runscope_bucket",
		F:    testAccSweepBuckets,
	})
}

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
				Config: fmt.Sprintf(testAccRunscopeBucketBasicConfig, bucketName, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists("runscope_bucket.bucket", &bucket),
					resource.TestCheckResourceAttr("runscope_bucket.bucket", "name", bucketName),
					resource.TestCheckResourceAttr("runscope_bucket.bucket", "team_uuid", teamId),
					resource.TestCheckResourceAttrSet("runscope_bucket.bucket", "default"),
					resource.TestCheckResourceAttrSet("runscope_bucket.bucket", "verify_ssl"),
					resource.TestCheckResourceAttrSet("runscope_bucket.bucket", "trigger_url"),
				),
			},
			{
				ResourceName:      "runscope_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckBucketDestroy(s *terraform.State) error {
	ctx := context.Background()
	client := testAccProvider.Meta().(*providerConfig).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "runscope_bucket" {
			if _, err := client.Bucket.Get(ctx, &runscope.BucketGetOpts{Key: rs.Primary.ID}); err == nil {
				return fmt.Errorf("Record %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckBucketExists(n string, b *runscope.Bucket) resource.TestCheckFunc {
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

		bucket, err := client.Bucket.Get(ctx, &runscope.BucketGetOpts{Key: rs.Primary.ID})
		if err != nil {
			return err
		}

		if bucket.Key != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*b = *bucket

		return nil
	}
}

const testAccRunscopeBucketBasicConfig = `
resource "runscope_bucket" "bucket" {
  name      = "%s"
  team_uuid = "%s"
}`

func testAccSweepBuckets(_ string) error {
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
}
