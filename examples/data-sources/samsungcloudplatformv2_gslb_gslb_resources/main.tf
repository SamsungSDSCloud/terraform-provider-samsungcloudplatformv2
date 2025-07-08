provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_gslb_gslb_resources" "gslbresources" {
  gslb_id = var.gslb_id
  size = var.size
}
