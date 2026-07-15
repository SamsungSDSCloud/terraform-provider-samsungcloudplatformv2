output "success_ids" {
  description = "Successfully removed accounts with IDs and names"
  value       = samsungcloudplatformv2_organization_account_remove.remove[*].success_ids
}

output "failed_ids" {
  description = "Failed accounts with error details"
  value       = samsungcloudplatformv2_organization_account_remove.remove[*].failed_ids
}