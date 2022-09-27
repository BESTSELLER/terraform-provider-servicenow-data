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
  sn_api_url  = "..."
  sn_api_user = "..."
  sn_api_pass = "..."
}
