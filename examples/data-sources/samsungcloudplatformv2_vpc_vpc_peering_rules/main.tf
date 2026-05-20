provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_vpc_peering_rules" "my_vpc_peering_rule_list" {
  vpc_peering_id       = var.vpc_peering_id
  size                 = var.size
  page                 = var.page
  sort                 = var.sort
  id                   = var.id
  source_vpc_id        = var.source_vpc_id
  source_vpc_type      = var.source_vpc_type
  destination_vpc_id   = var.destination_vpc_id
  destination_vpc_type = var.destination_vpc_type
  destination_cidr     = var.destination_cidr
  state                = var.state
}
