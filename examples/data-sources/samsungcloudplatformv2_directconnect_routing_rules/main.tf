provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_directconnect_routing_rules" "routingrules" {
  limit = var.limit
  direct_connect_id = var.direct_connect_id
}
