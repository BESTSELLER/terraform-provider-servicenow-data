package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/client"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/models"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ServiceCatalogOrderResourceName = "servicenow-data_servicecatalog_order"

func ServiceCatalogOrderResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: serviceCatalogOrderCreate,
		ReadContext:   serviceCatalogOrderRead,
		// UpdateContext:  serviceCatalogOrderUpdate, SNOW catalog orders are immutable
		DeleteContext: serviceCatalogOrderDelete,
		Schema: *models.MergeSchema(
			models.ServiceCatalogOrderColumns,
			*models.MergeSchema(
				models.ServiceCatalogOrderRequestColumns,
				models.ServiceCatalogOrderResponseColumns,
			),
		),
		SchemaVersion:  1,
		StateUpgraders: nil,
		CustomizeDiff:  nil,
		Importer: &schema.ResourceImporter{
			StateContext: nil,
		},
		DeprecationMessage: "",
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(5 * time.Second),
			Read:    schema.DefaultTimeout(5 * time.Second),
			Update:  schema.DefaultTimeout(5 * time.Second),
			Delete:  schema.DefaultTimeout(5 * time.Second),
			Default: schema.DefaultTimeout(5 * time.Second),
		},
		Description: `
A Service Now ServiceCatalog order. It consumes the "order_now" API ( https://developer.servicenow.com/dev.do#!/reference/api/utah/rest/c_ServiceCatalogAPI#servicecat-POST-items-order_now )

Please note that Service Now catalog orders are immutable.
That means if you try to destroy an order, it will be removed from terraform but nothing will happen in Service Now.
Also, if you try to update an order, it will be removed from terraform and a new one will be created. In other words, you will end up with two orders in service now.
`,
		UseJSONNumber: false,
	}
}

func prepareCreatePayload(d *schema.ResourceData) map[string]interface{} {
	payload := map[string]interface{}{}
	for k, _ := range models.ServiceCatalogOrderRequestColumns {
		payload[k] = d.Get(k)
	}
	return payload
}

func sendCreatePayload(ctx context.Context, url string, payload map[string]interface{}, client *client.Client) (map[string]any, diag.Diagnostics) {
	rawData, err := client.SendRequest(http.MethodPost, url, payload, 200)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("Error on SendRequest: %s", err))
	}

	var data map[string]any
	err = json.Unmarshal(*rawData, &data)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("Error on json.Unmarsha: %s", err))
	}
	return data, nil
}

func saveCreatePayload(data_result map[string]any, data *schema.ResourceData) diag.Diagnostics {
	for k, v := range data_result {
		if k == models.ServiceCatalogOrderBlacklistedField {
			continue
		}
		if err := data.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func serviceCatalogOrderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	url := fmt.Sprintf("/api/sn_sc/v1/servicecatalog/items/%s/order_now", d.Get("sc_cat_item_id"))
	payload := prepareCreatePayload(d)
	tflog.Info(ctx, fmt.Sprintf("serviceCatalogOrderCreate: url=%s, params=%s", url, payload))

	client := m.(*client.Client)
	data, diags := sendCreatePayload(ctx, url, payload, client)
	if diags != nil {
		return diags
	}

	tflog.Info(ctx, fmt.Sprintf("serviceCatalogOrderCreate: result=%s", data))
	data_result := data["result"].(map[string]any)

	diags = saveCreatePayload(data_result, d)
	if diags != nil {
		return diags
	}

	sys_id := data_result["sys_id"].(string)

	d.SetId(sys_id)
	return nil
}

func serviceCatalogOrderRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	sys_id := data.Id()
	tflog.Info(ctx, fmt.Sprintf("serviceCatalogOrderRead: reading a ServiceCatalogOrder in service now is not possible. This is a NO OP by design. Use servicenow-data_table_row with the appropriate table instead. sys_id=%s", sys_id))
	return nil
}

func serviceCatalogOrderDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	sys_id := data.Id()
	tflog.Info(ctx, fmt.Sprintf("serviceCatalogOrderDelete: deleting in service now is not possible. This is a NO OP by design. The resource will be removed from the terraform state. sys_id=%s", sys_id))
	return nil
}
