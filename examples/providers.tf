terraform {
  required_providers {
    servicenow-data = {
      source  = "local/providers/servicenow-data"
      version = "1.0.0"
    }
    # add other providers here
  }
  required_version = ">= 0.13"
}
provider "servicenow-data" {
#  sn_api_url  = "..."
#  sn_api_user = "..."
#  sn_api_pass = "..."
}

data "servicenow-data_table_row" "test" {
  table_id = "x_beas_team_engi_0_approval_items"
  sys_id   = "82dc8b2dc3029910a1ec2a4ce0013134"

}

resource "local_file" "remote_state" {
  content  = jsonencode( data.servicenow-data_table_row.test)
  filename = "data_out.json"
}