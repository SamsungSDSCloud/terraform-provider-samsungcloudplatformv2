output "count" {
  description = "Total number of targets bound to the guardrail"
  value       = data.samsungcloudplatformv2_cloudcontrol_guardrail_bindings_targets.targets.total_count
}
output "sort" {
  description = "Sort criteria actually applied by the server"
  value       = data.samsungcloudplatformv2_cloudcontrol_guardrail_bindings_targets.targets.sort_result
}
output "size" {
  description = "Page size actually applied by the server"
  value       = data.samsungcloudplatformv2_cloudcontrol_guardrail_bindings_targets.targets.size
}
output "page" {
  description = "Page number actually applied by the server"
  value       = data.samsungcloudplatformv2_cloudcontrol_guardrail_bindings_targets.targets.page
}
output "targets" {
  description = "List of targets bound to the guardrail (email only populated when target_type=ACCOUNT)"
  value       = data.samsungcloudplatformv2_cloudcontrol_guardrail_bindings_targets.targets.targets
}