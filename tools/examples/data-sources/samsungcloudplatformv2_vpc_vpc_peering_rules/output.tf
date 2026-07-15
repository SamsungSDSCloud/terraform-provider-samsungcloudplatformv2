output "vpc_peering_rules" {
  value = {
    count : data.samsungcloudplatformv2_vpc_vpc_peering_rules.my_vpc_peering_rule_list.total_count,
    page : data.samsungcloudplatformv2_vpc_vpc_peering_rules.my_vpc_peering_rule_list.page,
    size : data.samsungcloudplatformv2_vpc_vpc_peering_rules.my_vpc_peering_rule_list.size,
    sort : data.samsungcloudplatformv2_vpc_vpc_peering_rules.my_vpc_peering_rule_list.sort_final,
    vpc_peering_rules : data.samsungcloudplatformv2_vpc_vpc_peering_rules.my_vpc_peering_rule_list.vpc_peering_rules
  }
}
