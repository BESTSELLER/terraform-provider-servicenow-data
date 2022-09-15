package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ServiceNowDataProvider() *schema.Provider {
	dataSources := make(map[string]*schema.Resource)
	resources := make(map[string]*schema.Resource)
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"sn_api_user": {
				Type:        schema.TypeString,
				Optional:    false,
				DefaultFunc: schema.EnvDefaultFunc("SN_API_USER", ""),
				Description: "The user required to auth to the SN table API using basic auth"},
			"sn_api_pass": {
				Type:        schema.TypeString,
				Optional:    false,
				DefaultFunc: schema.EnvDefaultFunc("SN_API_PASS", ""),
				Description: "The Password required to auth to the SN table API using basic auth"},
			"sn_api_url": {
				Type:        schema.TypeString,
				Optional:    false,
				DefaultFunc: schema.EnvDefaultFunc("SN_API_URL", ""),
				Description: "The URL to the SN table using basic auth"},
		},

		ResourcesMap:   resources,
		DataSourcesMap: dataSources,
	}
	return p
}
