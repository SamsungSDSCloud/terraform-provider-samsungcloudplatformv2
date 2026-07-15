provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_transit_gateway_firewall" "my_tgw_firewall" {
  transit_gateway_id = var.transit_gateway_id
  product_type       = var.product_type
}

