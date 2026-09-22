output "count" {
  description = "Total count of targets assigned to the policy"
  value       = data.samsungcloudplatformv2_organization_targets_for_policy.targets.total_count
}

output "page" {
  description = "Current page number"
  value       = data.samsungcloudplatformv2_organization_targets_for_policy.targets.page
}

output "size" {
  description = "Page size"
  value       = data.samsungcloudplatformv2_organization_targets_for_policy.targets.size
}

output "sort" {
  description = "Sort criteria from response"
  value       = data.samsungcloudplatformv2_organization_targets_for_policy.targets.sort_result
}

output "targets" {
  description = "List of targets assigned to the policy"
  value       = data.samsungcloudplatformv2_organization_targets_for_policy.targets.targets
}
