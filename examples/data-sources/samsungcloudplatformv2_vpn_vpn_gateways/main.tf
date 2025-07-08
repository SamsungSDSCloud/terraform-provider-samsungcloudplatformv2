provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpn_vpn_gateways" "ids" {
  size = var.size
  page = var.page
  sort = var.sort
  name = var.name
  ip_address = var.ip_address
  vpc_id = var.vpc_id
  vpc_name = var.vpc_name
}
