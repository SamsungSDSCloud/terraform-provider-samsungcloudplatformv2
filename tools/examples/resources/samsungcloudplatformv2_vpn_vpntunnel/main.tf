provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpn_vpn_tunnel" "vpntunnel" {
  description = var.description
  name = var.name
  phase1 = var.phase1
  phase2 = var.phase2
  tags = var.tags
  vpn_gateway_id = var.vpn_gateway_id
}

