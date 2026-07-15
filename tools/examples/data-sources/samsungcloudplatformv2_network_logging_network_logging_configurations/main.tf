provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_network_logging_network_logging_configurations" "networkloggingconfigurations" {
  limit = var.limit
  resource_id = var.resource_id
  resource_name = var.resource_name
  resource_type = var.resource_type
}
