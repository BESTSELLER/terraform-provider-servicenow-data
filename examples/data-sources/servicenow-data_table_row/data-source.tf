data "servicenow-data_table_row" "example-with-sys_id" {
  table_id = "sys_user"
  sys_id   = "7a9dde3e6fa4310005a9fbf7eb3ee495"
}

data "servicenow-data_table_row" "example-with-email-query" {
  table_id = "sys_user"
  row_data = {
    "email" : "example.value@example.com"
  }
}