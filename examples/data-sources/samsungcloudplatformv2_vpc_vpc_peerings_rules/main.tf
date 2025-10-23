provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_vpc_peering_rules" "my_vpc_peering_rule_list" {
  vpc_peering_id = var.vpc_peering_id
  size           = var.size
}
