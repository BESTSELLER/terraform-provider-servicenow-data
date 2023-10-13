package provider

import (
	"context"

	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/client"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/datasource"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ServiceNowDataProvider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"sn_api_user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SN_API_USER", ""),
				Description: "The user required to auth to the SN table API using basic auth"},
			"sn_api_pass": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SN_API_PASS", ""),
				Description: "The Password required to auth to the SN table API using basic auth"},
			"sn_api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SN_API_URL", ""),
				Description: "The URL to the SN table using basic auth"},
		},

		ResourcesMap: map[string]*schema.Resource{
			resource.TableRowResourceName:            resource.TableRowResource(),
			resource.ServiceCatalogOrderResourceName: resource.ServiceCatalogOrderResource(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			resource.TableRowResourceName: datasource.TableRowDatasource(),
		},
		ConfigureContextFunc: providerConfigure,
	}
	return p
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var host, username, password *string
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	hVal, ok := d.GetOk("sn_api_user")
	if ok {
		tempHost := hVal.(string)
		username = &tempHost
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "sn_api_user - Missing required config",
		})
	}
	hVal, ok = d.GetOk("sn_api_pass")
	if ok {
		tempHost := hVal.(string)
		password = &tempHost
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "sn_api_pass - Missing required config",
		})
	}
	hVal, ok = d.GetOk("sn_api_url")
	if ok {
		tempHost := hVal.(string)
		host = &tempHost
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "sn_api_url - Missing required config",
		})
	}
	if diags.HasError() {
		return nil, diags
	}

	c := client.NewClient(ctx, *host, *username, *password)

	return c, diags
}
