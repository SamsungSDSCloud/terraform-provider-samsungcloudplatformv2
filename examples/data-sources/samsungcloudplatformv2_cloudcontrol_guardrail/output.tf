output "id" {
  value = data.samsungcloudplatformv2_cloudcontrol_guardrail.guardrail.guardrail.id
}

output "name" {
  value = data.samsungcloudplatformv2_cloudcontrol_guardrail.guardrail.guardrail.name
}

output "description" {
  value = data.samsungcloudplatformv2_cloudcontrol_guardrail.guardrail.guardrail.description
}

output "guidance" {
  value = data.samsungcloudplatformv2_cloudcontrol_guardrail.guardrail.guardrail.guidance
}

output "service_name" {
  value = data.samsungcloudplatformv2_cloudcontrol_guardrail.guardrail.guardrail.service_name
}

output "status" {
  value = data.samsungcloudplatformv2_cloudcontrol_guardrail.guardrail.guardrail.status
}

output "type" {
  value = data.samsungcloudplatformv2_cloudcontrol_guardrail.guardrail.guardrail.type
}

output "created_at" {
  value = data.samsungcloudplatformv2_cloudcontrol_guardrail.guardrail.guardrail.created_at
}

output "created_by" {
  value = data.samsungcloudplatformv2_cloudcontrol_guardrail.guardrail.guardrail.created_by
}

output "modified_at" {
  value = data.samsungcloudplatformv2_cloudcontrol_guardrail.guardrail.guardrail.modified_at
}

output "modified_by" {
  value = data.samsungcloudplatformv2_cloudcontrol_guardrail.guardrail.guardrail.modified_by
}

output "srn" {
  value = data.samsungcloudplatformv2_cloudcontrol_guardrail.guardrail.guardrail.srn
}

output "binding_ous" {
  description = "Organization Units this guardrail is bound to"
  value       = data.samsungcloudplatformv2_cloudcontrol_guardrail.guardrail.guardrail.binding_ous
}