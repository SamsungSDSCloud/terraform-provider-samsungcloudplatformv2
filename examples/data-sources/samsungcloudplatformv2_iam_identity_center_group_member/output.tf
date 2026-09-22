output "group_member" {
  description = "IAM Identity Center Group Member"
  value       = data.samsungcloudplatformv2_iam_identity_center_group_member.example
}

output "id" {
  description = "Group Member ID"
  value       = data.samsungcloudplatformv2_iam_identity_center_group_member.example.id
}

output "user_uuids" {
  description = "List of User IDs in the group"
  value       = data.samsungcloudplatformv2_iam_identity_center_group_member.example.user_uuids
}
