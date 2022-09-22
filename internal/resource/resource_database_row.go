package resource

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

const DatabaseRowResourceName = "servicenow-data_table_row"

func DatabaseRowResource() *schema.Resource {
	return &schema.Resource{
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
