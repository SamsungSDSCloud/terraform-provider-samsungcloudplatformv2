provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpn_vpn_tunnel" "vpntunnel" {
  name = var.name
  vpn_gateway_id = var.vpn_gateway_id
  phase1 = var.phase1
  phase2 = var.phase2
  description = var.description
  tags = var.tags
}

