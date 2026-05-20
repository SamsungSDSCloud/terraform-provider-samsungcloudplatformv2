provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_transit_gateway_uplink_rule" "uplink" {
  description        = var.description
  destination_cidr   = var.destination_cidr
  destination_type   = var.destination_type
  transit_gateway_id = var.transit_gateway_id
}
