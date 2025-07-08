provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_network_logging_network_logging_storages" "networkloggingstorages" {
  resource_type = var.resource_type
}
