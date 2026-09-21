output "group_member" {
  description = "IAM Identity Center Group Member"
  value       = samsungcloudplatformv2_iam_identity_center_group_member.example
}

output "id" {
  description = "Group Member ID"
  value       = samsungcloudplatformv2_iam_identity_center_group_member.example.id
}

output "member_id" {
  description = "ID of the user added as a group member"
  value       = samsungcloudplatformv2_iam_identity_center_group_member.example.member_id
}
