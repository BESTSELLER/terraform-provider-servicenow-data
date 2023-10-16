resource "servicenow-data_servicecatalog_order" "order" {
  sc_cat_item_id = "e2a2e5bc1b757850d5a68773604bcb32"
  variables = {
    "var1" : "value1",
    "var2" : "value2",
    "var3" : "value3",
  }
}

# sample output
# "result": {
#     "$$uiNotification": [],
#     "number": "REQ0227772",
#     "parent_id": null,
#     "parent_table": "task",
#     "request_id": "24e3e99e97f53110af6574971153afaa",
#     "request_number": "REQ0227772",
#     "sys_id": "24e3e99e97f53110af6574971153afaa",
#     "table": "sc_request"
# }

data "servicenow-data_table_row" "request" {
  table_id = "sc_request"
  row_data = {
    "number" : servicenow-data_servicecatalog_order.order.request_number
  }
}

data "servicenow-data_table_row" "request_item" {
  table_id = "sc_req_item"
  row_data = {
    "request" : data.servicenow-data_table_row.request.sys_id
  }
}
