provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpn_vpn_gateway" "vpngateway" {
  id = var.id
}
