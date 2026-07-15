provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_gslb_gslb" "gslb" {
  gslb_create = var.gslb
  tags        = var.tag
}