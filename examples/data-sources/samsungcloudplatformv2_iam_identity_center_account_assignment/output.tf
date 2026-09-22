output "principal_id" {
  description = "Principal ID"
  value       = data.samsungcloudplatformv2_iam_identity_center_account_assignment.example.principal_id
}

output "principal_name" {
  description = "Principal Name"
  value       = data.samsungcloudplatformv2_iam_identity_center_account_assignment.example.principal_name
}

output "principal_type" {
  description = "Principal Type"
  value       = data.samsungcloudplatformv2_iam_identity_center_account_assignment.example.principal_type
}

output "target_account_name" {
  description = "Target Account Name"
  value       = data.samsungcloudplatformv2_iam_identity_center_account_assignment.example.target_account_name
}