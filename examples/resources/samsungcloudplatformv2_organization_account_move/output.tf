output "success_ids" {
  description = "Successfully moved account information"
  value       = samsungcloudplatformv2_organization_account_move.move[*].success_ids
}