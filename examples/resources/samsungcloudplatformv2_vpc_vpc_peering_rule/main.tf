provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_vpc_peering_rule" "create_peering_rule" {
  vpc_peering_id        = var.vpc_peering_id
  destination_cidr      = var.destination_cidr
  destination_vpc_type  = var.destination_vpc_type
  tags                  = var.tags
}
