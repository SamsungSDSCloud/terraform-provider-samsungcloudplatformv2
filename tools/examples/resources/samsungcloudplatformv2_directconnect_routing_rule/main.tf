provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_directconnect_routing_rule" "routingrule" {
  direct_connect_id = var.routing_rule_direct_connect_id
  destination_type = var.routing_rule_destination_type
  destination_cidr = var.routing_rule_destination_cidr
}
