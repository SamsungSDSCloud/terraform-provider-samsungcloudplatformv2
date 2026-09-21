output "user" {
  description = "IAM Identity Center User"
  value       = samsungcloudplatformv2_iam_identity_center_user.example
  sensitive   = true
}

output "user_id" {
  description = "User ID"
  value       = samsungcloudplatformv2_iam_identity_center_user.example.id
}

output "user_name" {
  description = "User Name"
  value       = samsungcloudplatformv2_iam_identity_center_user.example.name
}