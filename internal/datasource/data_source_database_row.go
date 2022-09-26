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

// TODO
func DatabaseRowDatasource() *schema.Resource {
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
		ReadContext:   DatabaseRowRead,
		Description:   "A row in a SN table",
		UseJSONNumber: false,
	}
}

func DatabaseRowRead(_ context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var tableID, rowID string
	var err error
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	if data.Id() != "" {
		tableID, rowID, err = ExtractIDs(data)
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
		rowID = data.Get("sys_id").(string)
		if rowID == "" {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "sys_id is mandatory",
			})
		}
	}

	rowData, err := c.GetTableRow(tableID, rowID)
	if err != nil {
		return diag.FromErr(err)
	}

	diags = append(diags, ParsedResultToSchema(data, rowData)...)

	data.SetId(fmt.Sprintf("%s/%s", tableID, rowID))

	return diags
}

func ExtractIDs(data *schema.ResourceData) (tableID, rowID string, err error) {
	ids := strings.Split(data.Id(), `/`)
	if len(ids) != 2 {
		return "", "", errors.New(fmt.Sprintf("Faulty id!%s", data.Id()))
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
