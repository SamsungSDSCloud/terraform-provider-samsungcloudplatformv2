provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpn_vpn_tunnel" "vpntunnel" {
  id = var.id
}
