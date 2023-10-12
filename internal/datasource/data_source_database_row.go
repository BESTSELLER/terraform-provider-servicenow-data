package datasource

import (
	"context"
	"fmt"

	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/client"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/models"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TableRowDatasource() *schema.Resource {
	return &schema.Resource{
		Schema: *models.MergeSchema(models.DefaultSystemColumns, map[string]*schema.Schema{
			"sys_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"row_data": {
				Description: "Columns",
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString},
			},
		}),
		ReadContext:   tableRowRead,
		Description:   "A row in a SN table",
		UseJSONNumber: false,
	}
}

func tableRowRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("tableRowRead: data=%+v", data))
	c := m.(*client.Client)
	var err error
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var payload map[string]interface{}
	tableID, ok := data.GetOk("table_id")
	if !ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "tableID is mandatory",
		})
		return diags
	}
	sysID, ok := data.GetOk("sys_id")
	if !ok {
		payload = data.Get("row_data").(map[string]interface{})
		if payload == nil || len(payload) == 0 {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "sys_id or row_data is mandatory",
			})
		}
	} else {
		payload = map[string]interface{}{"sys_id": sysID}
	}
	rowData, err := c.GetTableRow(tableID.(string), payload)
	if err != nil {
		return diag.FromErr(err)
	}

	if _, exists := (*rowData).SysData["sys_id"]; !exists {
		return append(diags, diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("No row found in ServiceNow for %+v - %+v", data.Get("sys_id"), data.Get("row_data")),
		}}...)
	}
	diags = append(diags, resource.ParsedResultToSchema(data, rowData)...)

	data.SetId(fmt.Sprintf("%s/%s", tableID, rowData.SysData["sys_id"]))
	return diags
}
