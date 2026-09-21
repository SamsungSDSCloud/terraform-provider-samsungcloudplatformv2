output "permission_set_policies" {
  description = "IAM Identity Center Permission Set Policy"
  value       = data.samsungcloudplatformv2_iam_identity_center_permission_set_policies.example
}

output "id" {
  description = "Permission Set Policy ID"
  value       = data.samsungcloudplatformv2_iam_identity_center_permission_set_policies.example.id
}
