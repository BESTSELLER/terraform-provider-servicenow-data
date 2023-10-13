package models

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DefaultSystemColumns = map[string]*schema.Schema{
	"table_id": {
		Type:     schema.TypeString,
		Required: true,
	},
	"sys_created_by": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"sys_created_on": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"sys_mod_count": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"sys_tags": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"sys_updated_by": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"sys_updated_on": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"sys_domain": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"sys_domain_path": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"sys_name": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"sys_package": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"sys_policy": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"sys_update_name": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"sys_class_name": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	}}

type ReferenceItem struct {
	Link  string `json:"link"`
	Value string `json:"value"`
}

type RawResult struct {
	Result map[string]json.RawMessage `json:"result"`
}

type RawResultList struct {
	Result []map[string]json.RawMessage `json:"result"`
}

type ParsedResult struct {
	SysData map[string]string `json:"sys_data"`
	RowData map[string]string `json:"row_data"`
}

func MergeSchema(data1, data2 map[string]*schema.Schema) *map[string]*schema.Schema {
	result := make(map[string]*schema.Schema, len(data1)+len(data2))
	for k, v := range data1 {
		result[k] = v
	}
	for k, v := range data2 {
		result[k] = v
	}
	return &result
}
