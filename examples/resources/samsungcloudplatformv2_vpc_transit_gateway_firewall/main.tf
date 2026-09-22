provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_transit_gateway_firewall" "my_tgw_firewall" {
  transit_gateway_id  = var.transit_gateway_id
  product_type        = var.product_type
  uplink_active_zone  = var.uplink_active_zone
  uplink_standby_zone = var.uplink_standby_zone
}

