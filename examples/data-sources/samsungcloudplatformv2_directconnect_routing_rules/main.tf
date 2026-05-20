provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_directconnect_routing_rules" "routingrules" {
  size              = var.size
  page              = var.page
  sort              = var.sort
  id                = var.id
  direct_connect_id = var.direct_connect_id
  destination_type  = var.destination_type
  destination_cidr  = var.destination_cidr
  state             = var.state
}
