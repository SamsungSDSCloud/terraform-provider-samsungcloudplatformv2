provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_network_logging_network_logging_configurations" "networkloggingconfigurations" {
  resource_type = var.resource_type
}
