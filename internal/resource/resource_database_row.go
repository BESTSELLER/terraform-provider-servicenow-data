package resource

import (
	"context"
	"fmt"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/client"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/datasource"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
	"time"
)

const TableRowResourceName = "servicenow-data_table_row"

func TableRowResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: tableRowCreate,
		ReadContext:   datasource.TableRowRead,
		UpdateContext: tableRowUpdate,
		DeleteContext: tableRowDelete,
		Schema: *models.MergeSchema(models.DefaultSystemColumns, map[string]*schema.Schema{
			"sys_id": {
				Type:     schema.TypeString,
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
		Description:   "A row in a SN table",
		UseJSONNumber: false,
	}
}

func tableRowCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	tableID := d.Get("table_id").(string)
	tableData := d.Get("row_data")
	insertResult, err := c.InsertTableRow(tableID, tableData)

	if err != nil {
		return diag.FromErr(err)
	}

	diags = append(diags, datasource.ParsedResultToSchema(d, insertResult)...)
	sysID, ok := d.GetOk("sys_id")
	if !ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "sys_id is mandatory, row create did not return a sys_id",
		})
		return diags
	}
	err = d.Set("sys_id", sysID)
	d.SetId(fmt.Sprintf("%s/%s", tableID, sysID))
	diags = append(diags, diag.Diagnostic{Severity: diag.Warning,
		Summary: d.Id()})
	return diags
}

func tableRowUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var tableID, sysID string
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if d.HasChange("table_id") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "table_id cannot be modified after creation",
		})
		return diags
	}
	if d.HasChange("sys_id") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "sys_id cannot be modified after creation",
		})
		return diags
	}

	tableID, sysID, err := datasource.ExtractIDs(d)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	tableData := d.Get("row_data")
	rowData, err := c.UpdateTableRow(tableID, sysID, tableData)
	if err != nil {
		return diag.FromErr(err)
	}

	diags = append(diags, datasource.ParsedResultToSchema(d, rowData)...)

	return diags
}

func tableRowDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	split := strings.Split(d.Id(), `\`)
	tableID := split[0]
	sysID := split[1]

	err := c.DeleteTableRow(tableID, sysID)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	return diags
}
