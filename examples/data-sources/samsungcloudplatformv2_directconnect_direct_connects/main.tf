provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_directconnect_direct_connects" "directconnects" {
  limit = var.limit
}
