package resource

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

const DatabaseRowResourceName = "sn_application"

func DatabaseRowResource() *schema.Resource {
	return &schema.Resource{
		Schema:         RowSchema,
		SchemaVersion:  1,
		StateUpgraders: nil,
		CreateContext:  nil,
		ReadContext:    nil,
		UpdateContext:  nil,
		DeleteContext:  nil,
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

var RowSchema = map[string]*schema.Schema{
	"table_id": {
		Type:     schema.TypeString,
		Required: true,
		Computed: false,
	},
	"row_data": {
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"sys_id": {
					Description: "The unique id of the row",
					Type:        schema.TypeString,
					Required:    false,
					Computed:    true},
				"sys_updated_by": {
					Description: "User that made the last update",
					Type:        schema.TypeString,
					Required:    false,
					Computed:    true},
				"sys_created_by": {
					Description: "Account that created the row",
					Type:        schema.TypeString,
					Required:    false,
					Computed:    true},
				"sys_created_on": {
					Description: "Creation Time",
					Type:        schema.TypeString,
					Required:    false,
					Computed:    true},
				"sys_updated_on": {
					Description: "Last update Time",
					Type:        schema.TypeString,
					Required:    false,
					Computed:    true},
				"custom_columns": {
					Description: "Custom columns that are not references",
					Type:        schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString},
				},
			},
		},
	},
}
