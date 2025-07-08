provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_gslb_gslbs" "gslbs" {
  size = var.size
}
