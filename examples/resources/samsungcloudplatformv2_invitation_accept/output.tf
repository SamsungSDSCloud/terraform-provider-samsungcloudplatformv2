output "master_account_email" {
  description = "Master Account Email"
  value       = samsungcloudplatformv2_invitation_accept.accept.master_account_email
}

output "success_id" {
  description = "Success ID object"
  value       = samsungcloudplatformv2_invitation_accept.accept.success_id
}

output "failed_id" {
  description = "Failed ID object"
  value       = samsungcloudplatformv2_invitation_accept.accept.failed_id
}