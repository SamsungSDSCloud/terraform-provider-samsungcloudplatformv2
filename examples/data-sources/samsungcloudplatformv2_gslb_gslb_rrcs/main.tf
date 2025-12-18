provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_gslb_gslb_rrcs" "my_gslb_rrc_list" {
  size   = var.size
  page   = var.page
  sort   = var.sort
  region = var.region
  status = var.status
  name   = var.name
}
