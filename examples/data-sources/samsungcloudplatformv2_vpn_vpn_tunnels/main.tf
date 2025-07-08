provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpn_vpn_tunnels" "ids" {
  size = var.size
  page = var.page
  sort = var.sort
  name = var.name
  vpn_gateway_id = var.vpn_gateway_id
  vpn_gateway_name = var.vpn_gateway_name
  peer_gateway_ip = var.peer_gateway_ip
  remote_subnet = var.remote_subnet
}
