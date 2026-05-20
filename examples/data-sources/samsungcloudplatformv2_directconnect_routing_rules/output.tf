output "routingRules" {
  value = {
    count: data.samsungcloudplatformv2_directconnect_routing_rules.routingrules.total_count,
    routing_rules: data.samsungcloudplatformv2_directconnect_routing_rules.routingrules.routing_rules,
    page: data.samsungcloudplatformv2_directconnect_routing_rules.routingrules.page,
    size: data.samsungcloudplatformv2_directconnect_routing_rules.routingrules.size,
    sort: data.samsungcloudplatformv2_directconnect_routing_rules.routingrules.sort_final,
  }
}