provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_gslb_gslb_rrc_update" "my_gslb_rrc" {
  gslb_id = var.gslb_id
  region  = var.region
  status  = var.status
}
