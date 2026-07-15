output "vpcpeeringrule_output" {
  value = {
    vpc_peering_rule: samsungcloudplatformv2_vpc_vpc_peering_rule.create_peering_rule.vpc_peering_rule
  }
}
