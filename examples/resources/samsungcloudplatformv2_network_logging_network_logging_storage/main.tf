provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_network_logging_network_logging_storage" "networkloggingstorage" {
  bucket_name = var.bucket_name
  resource_type = var.resource_type
}
