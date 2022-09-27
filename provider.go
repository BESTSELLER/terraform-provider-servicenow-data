package main

import (
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return provider.ServiceNowDataProvider()
}
