package resource

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

const databaseRowResourceName = "sn_application"

// TODO
func databaseRowResource() *schema.Resource {
	return &schema.Resource{
		Schema:         nil,
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
