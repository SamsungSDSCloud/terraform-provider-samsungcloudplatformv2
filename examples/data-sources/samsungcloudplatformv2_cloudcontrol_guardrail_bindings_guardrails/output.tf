output "count" {
  description = "Total number of guardrails bound to the target"
  value       = data.samsungcloudplatformv2_cloudcontrol_guardrail_bindings_guardrails.guardrails.total_count
}
output "sort" {
  description = "Sort criteria actually applied by the server"
  value       = data.samsungcloudplatformv2_cloudcontrol_guardrail_bindings_guardrails.guardrails.sort_result
}
output "size" {
  description = "Page size actually applied by the server"
  value       = data.samsungcloudplatformv2_cloudcontrol_guardrail_bindings_guardrails.guardrails.size
}
output "page" {
  description = "Page number actually applied by the server"
  value       = data.samsungcloudplatformv2_cloudcontrol_guardrail_bindings_guardrails.guardrails.page
}
output "guardrails" {
  description = "List of guardrails bound to the target (includes link_types: directed/inherited)"
  value       = data.samsungcloudplatformv2_cloudcontrol_guardrail_bindings_guardrails.guardrails.guardrails
}