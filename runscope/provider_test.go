package runscope

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		// provider is called terraform-provider-runscope ie runscope
		"runscope": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}

}

func TestProviderImpl(t *testing.T) {
	var _ = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("RUNSCOPE_ACCESS_TOKEN"); v == "" {
		t.Fatal("RUNSCOPE_ACCESS_TOKEN must be set for acceptance tests")
	}

	if v := os.Getenv("RUNSCOPE_TEAM_ID"); v == "" {
		t.Fatal("RUNSCOPE_TEAM_ID must be set for acceptance tests")
	}
}
