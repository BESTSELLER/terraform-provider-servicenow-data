package models

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var ServiceCatalogOrderBlacklistedField = "$$uiNotification"

var ServiceCatalogOrderColumns = map[string]*schema.Schema{
	// https://developer.servicenow.com/dev.do#!/reference/api/utah/rest/c_ServiceCatalogAPI#servicecat-PUT-items-submit_guide?navFilter=servicecatalog/items
	"sc_cat_item_id": {
		Description:  "The id of the catalog item to order",
		Type:         schema.TypeString,
		ValidateFunc: validation.StringLenBetween(32, 32),
		Required:     true,
		ForceNew:     true,
	},
}

var ServiceCatalogOrderRequestColumns = map[string]*schema.Schema{
	// REQUEST
	"sysparm_also_request_for": {
		Description: "Comma-separated string of user sys_ids of other users for which to order the specified item. User sys_ids are located in the User [sys_user] table.",
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
	},
	"sysparm_quantity": {
		Description: "Quantity of the item. Cannot be a negative number.",
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     1,
		ForceNew:    true,
	},
	"sysparm_requested_for": {
		Description:  "Sys_id of the user for whom to order the specified item. Located in the User [sys_user] table.",
		Type:         schema.TypeString,
		ValidateFunc: validation.StringLenBetween(32, 32),
		Optional:     true,
		ForceNew:     true,
	},
	"variables": {
		Description: "Name-value pairs of all mandatory cart item variables. Mandatory variables are defined on the associated form.",
		Required:    true,
		Type:        schema.TypeMap,
		ForceNew:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString},
	},
	"get_portal_messages": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "true",
		ForceNew: true,
	},
	"sysparm_no_validation": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "true",
		ForceNew: true,
	},
}

var ServiceCatalogOrderResponseColumns = map[string]*schema.Schema{
	// RESPONSE
	"sys_id": {
		Description: "Sys_id of the order.",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"number": {
		Description: "Number of the generated request",
		Type:        schema.TypeString,
		Computed:    true,
	},
	// this field is not documented, causes errors and is not needed
	// "$$uiNotification": {
	// 	Type:     schema.TypeSet,
	// 	Computed: true,
	// 	Elem:     &schema.Schema{Type: schema.TypeString},
	// 	Optional: true,
	// },
	"request_number": {
		Description: "Request number.",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"parent_id": {
		Description: "If available, the sys_id of the parent record from which the request is created",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"request_id": {
		Description: "Sys_id of the order request.",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"parent_table": {
		Description: "If available, the name of the parent table from which the request is created.",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"table": {
		Description: "Table name of the request.",
		Type:        schema.TypeString,
		Computed:    true,
	},
}
