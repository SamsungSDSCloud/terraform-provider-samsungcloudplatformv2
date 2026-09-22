output "user" {
  description = "IAM Identity Center User"
  value       = data.samsungcloudplatformv2_iam_identity_center_user.example
}

output "user_name" {
  description = "User Name"
  value       = data.samsungcloudplatformv2_iam_identity_center_user.example.name
}

output "user_email" {
  description = "User Email"
  value       = data.samsungcloudplatformv2_iam_identity_center_user.example.email
}