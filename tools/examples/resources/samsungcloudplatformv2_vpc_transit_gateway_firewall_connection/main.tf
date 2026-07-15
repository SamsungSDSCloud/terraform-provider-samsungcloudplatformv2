provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_transit_gateway_firewall_connection" "my_firewall_connection" {
  transit_gateway_id = var.transit_gateway_id
}
