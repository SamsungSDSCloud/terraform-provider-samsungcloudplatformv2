provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpn_vpn_gateway" "vpngateway" {
  ip_address = var.ip_address
  ip_id = var.ip_id
  ip_type = var.ip_type
  name = var.name
  vpc_id = var.vpc_id
  description = var.description
  tags = var.tags
}