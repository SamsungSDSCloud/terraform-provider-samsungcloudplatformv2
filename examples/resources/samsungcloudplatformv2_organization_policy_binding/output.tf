output "success_ids" {
  description = "List of IDs successfully bound"
  value       = samsungcloudplatformv2_organization_policy_binding.policy_binding.success_ids
}
output "failed_ids" {
  description = "List of IDs that failed to bind"
  value       = samsungcloudplatformv2_organization_policy_binding.policy_binding.failed_ids
}