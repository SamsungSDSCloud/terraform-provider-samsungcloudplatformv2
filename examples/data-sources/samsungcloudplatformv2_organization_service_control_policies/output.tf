output "scp_policies" {
  value = {
    policies    = data.samsungcloudplatformv2_organization_service_control_policies.scp_policies.policies
    total_count = data.samsungcloudplatformv2_organization_service_control_policies.scp_policies.total_count
    page        = data.samsungcloudplatformv2_organization_service_control_policies.scp_policies.page
    size        = data.samsungcloudplatformv2_organization_service_control_policies.scp_policies.size
    sort        = data.samsungcloudplatformv2_organization_service_control_policies.scp_policies.sort_result
  }
}