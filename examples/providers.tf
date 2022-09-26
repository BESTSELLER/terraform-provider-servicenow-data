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

data "servicenow-data_table_row" "LasseG" {
  table_id = "sys_user"
  sys_id   = "7a9dde3e6fa4310005a9fbf7eb3ee495"
}

data "servicenow-data_table_row" "AndreiP" {
  table_id = "sys_user"
  sys_id   = "254400dbdb0bc34032fe9ea9db96190b"
}

resource "servicenow-data_table_row" "eng-services-vault" {
  table_id = "x_beas_team_engi_0_lasse"
  row_data = {
    "team": "engineering-services2",
    "group_id_reader": "cd699222-ce5b-47ba-8d20-da254757c45c"
    "group_id_admin": "a8d94edc-8f08-4db7-a4c1-8e2a00d55795"
    "approvers": "${data.servicenow-data_table_row.LasseG.sys_id},${data.servicenow-data_table_row.AndreiP.sys_id}"
  }
}
