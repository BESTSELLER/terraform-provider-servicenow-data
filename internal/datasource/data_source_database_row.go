package datasource

import (
	"context"
	"errors"
	"fmt"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/client"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func TableRowDatasource() *schema.Resource {
	return &schema.Resource{
		Schema: *models.MergeSchema(models.DefaultSystemColumns, map[string]*schema.Schema{
			"sys_id": {
				Type:     schema.TypeString,
				Optional: true,
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
		ReadContext:   TableRowRead,
		Description:   "A row in a SN table",
		UseJSONNumber: false,
	}
}

func TableRowRead(_ context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var tableID, sysID string
	var err error
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	if data.Id() != "" {
		tableID, sysID, err = ExtractIDs(data.Id())
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	} else {
		tableID = data.Get("table_id").(string)
		if tableID == "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "tableID is mandatory",
			})
			return diags
		}
		sysID = data.Get("sys_id").(string)
		if sysID == "" {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "sys_id is mandatory",
			})
		}
	}

	rowData, err := c.GetTableRow(tableID, sysID)
	if err != nil {
		return diag.FromErr(err)
	}

	diags = append(diags, ParsedResultToSchema(data, rowData)...)

	data.SetId(fmt.Sprintf("%s/%s", tableID, sysID))

	return diags
}

func ExtractIDs(ID string) (tableID, sysID string, err error) {
	ids := strings.Split(ID, `/`)
	if len(ids) != 2 {
		return "", "", errors.New(fmt.Sprintf("Faulty id!%s", ID))
	}
	return ids[0], ids[1], nil
}

func ParsedResultToSchema(d *schema.ResourceData, result *models.ParsedResult) diag.Diagnostics {
	for k, v := range result.SysData {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("row_data", result.RowData); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
