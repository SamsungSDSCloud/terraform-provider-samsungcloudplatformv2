output "success_ids" {
  description = "Guardrail/unit pairs that were successfully bound"
  value       = samsungcloudplatformv2_cloudcontrol_guardrail_binding.guardrail_binding.success_ids
}
output "failed_ids" {
  description = "Guardrail/unit pairs that failed to bind"
  value       = samsungcloudplatformv2_cloudcontrol_guardrail_binding.guardrail_binding.failed_ids
}