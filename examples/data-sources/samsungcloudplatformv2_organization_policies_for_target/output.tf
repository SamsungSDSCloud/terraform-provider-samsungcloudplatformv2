output "count" {
  description = "Total count of policies assigned to the target"
  value       = data.samsungcloudplatformv2_organization_policies_for_target.policies.total_count
}

output "page" {
  description = "Current page number"
  value       = data.samsungcloudplatformv2_organization_policies_for_target.policies.page
}

output "size" {
  description = "Page size"
  value       = data.samsungcloudplatformv2_organization_policies_for_target.policies.size
}

output "sort" {
  description = "Sort criteria from response"
  value       = data.samsungcloudplatformv2_organization_policies_for_target.policies.sort_result
}

output "policies" {
  description = "List of policies assigned to the target"
  value       = data.samsungcloudplatformv2_organization_policies_for_target.policies.policies
}
