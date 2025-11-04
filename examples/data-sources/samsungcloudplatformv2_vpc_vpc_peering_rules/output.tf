output "vpc_peering_rules" {
  value = data.samsungcloudplatformv2_vpc_vpc_peering_rules.my_vpc_peering_rule_list.vpc_peering_rules
}