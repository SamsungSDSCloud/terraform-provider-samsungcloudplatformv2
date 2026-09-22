output "guardrails" {
  description = "List of guardrails matching the filter"
  value       = data.samsungcloudplatformv2_cloudcontrol_guardrails.guardrails.guardrails
}
output "count" {
  description = "Total number of guardrails matching the filter"
  value       = data.samsungcloudplatformv2_cloudcontrol_guardrails.guardrails.total_count
}
output "sort" {
  description = "Sort criteria actually applied by the server"
  value       = data.samsungcloudplatformv2_cloudcontrol_guardrails.guardrails.sort_result
}
output "size" {
  description = "Page size actually applied by the server"
  value       = data.samsungcloudplatformv2_cloudcontrol_guardrails.guardrails.size
}
output "page" {
  description = "Page number actually applied by the server"
  value       = data.samsungcloudplatformv2_cloudcontrol_guardrails.guardrails.page
}