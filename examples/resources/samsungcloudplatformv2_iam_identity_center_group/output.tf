output "group" {
  description = "IAM Identity Center Group"
  value       = samsungcloudplatformv2_iam_identity_center_group.example
}

output "group_id" {
  description = "Group ID"
  value       = samsungcloudplatformv2_iam_identity_center_group.example.id
}

output "group_name" {
  description = "Group Name"
  value       = samsungcloudplatformv2_iam_identity_center_group.example.group.name
}