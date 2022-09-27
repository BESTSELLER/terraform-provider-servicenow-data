package main

import (
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/provider"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = provider.ServiceNowDataProvider()
	testAccProviders = map[string]*schema.Provider{
		"sericenow-data": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}
func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("SN_API_URL"); err == "" {
		t.Fatal("SN_API_URL must be set for acceptance tests")
	}
	if err := os.Getenv("SN_API_USER"); err == "" {
		t.Fatal("SN_API_USER must be set for acceptance tests")
	}
	if err := os.Getenv("SN_API_PASS"); err == "" {
		t.Fatal("SN_API_PASS must be set for acceptance tests")
	}
}
