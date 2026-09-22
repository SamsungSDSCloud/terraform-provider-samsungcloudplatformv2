output "group" {
  description = "IAM Identity Center Group Data Source"
  value       = data.samsungcloudplatformv2_iam_identity_center_group.example
}

output "group_name" {
  description = "Group Name"
  value       = data.samsungcloudplatformv2_iam_identity_center_group.example.name
}

output "group_description" {
  description = "Group Description"
  value       = data.samsungcloudplatformv2_iam_identity_center_group.example.description
}