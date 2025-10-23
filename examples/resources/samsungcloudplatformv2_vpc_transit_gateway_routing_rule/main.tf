provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_transitgateway_routing_rule" "tgw_routingrule" {
  transit_gateway_id = var.routing_rule_transit_gateway_id
  description = var.routing_rule_description
  destination_type = var.routing_rule_destination_type
  destination_cidr = var.routing_rule_destination_cidr
  tgw_connection_vpc_id = var.routing_rule_tgw_connection_vpc_id
}
