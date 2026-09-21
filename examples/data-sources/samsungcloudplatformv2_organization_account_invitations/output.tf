output "count" {
  description = "Total Count"
  value       = data.samsungcloudplatformv2_organization_account_invitations.invitations.total_count
}

output "account_invitations" {
  description = "Account Invitations"
  value       = data.samsungcloudplatformv2_organization_account_invitations.invitations.account_invitations
}