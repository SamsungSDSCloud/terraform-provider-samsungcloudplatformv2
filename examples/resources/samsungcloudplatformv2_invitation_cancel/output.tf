output "success_ids" {
  description = "Successfully canceled invitation IDs"
  value       = samsungcloudplatformv2_invitation_cancel.cancel.canceled_ids
}