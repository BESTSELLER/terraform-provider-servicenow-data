package resource

import (
	"context"
	"fmt"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/client"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/datasource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
	"time"
)

const DatabaseRowResourceName = "servicenow-data_table_row"

func DatabaseRowResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: databaseRowCreate,
		ReadContext:   datasource.DatabaseRowRead,
		UpdateContext: databaseRowUpdate,
		DeleteContext: databaseRowDelete,
		Schema: map[string]*schema.Schema{
			"table_id": {
				Type:     schema.TypeString,
				Required: true,
			},
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
		},
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

func databaseRowCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	tableID := d.Get("table_id").(string)
	tableData := d.Get("row_data")
	insertResult, err := c.InsertTableRow(tableID, tableData)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("row_data", insertResult); err != nil {
		return diag.FromErr(err)
	}
	rowID, exists := (*insertResult)["sys_id"]
	if !exists {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "sys_id is mandatory, row create did not return a rowID",
		})
		return diags
	}
	err = d.Set("sys_id", rowID)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "failed to set sys_id",
		})
		return diags
	}

	d.SetId(fmt.Sprintf("%s/%s", tableID, rowID))

	return diags
}

func databaseRowUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	panic("NOT IMPLEMENTED!")
}

func databaseRowDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	split := strings.Split(d.Id(), `\`)
	tableID := split[0]
	rowID := split[1]

	err := c.DeleteTableRow(tableID, rowID)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	return diags
}
