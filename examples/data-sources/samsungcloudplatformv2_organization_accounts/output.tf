output "accounts" {
  description = "List of Organization Accounts"
  value       = data.samsungcloudplatformv2_organization_accounts.accounts.accounts
}

output "total_count" {
  description = "Total count of accounts"
  value       = data.samsungcloudplatformv2_organization_accounts.accounts.total_count
}

output "page" {
  description = "Current page number"
  value       = data.samsungcloudplatformv2_organization_accounts.accounts.page
}

output "size" {
  description = "Page size"
  value       = data.samsungcloudplatformv2_organization_accounts.accounts.size
}

output "sort" {
  description = "Sort criteria"
  value       = data.samsungcloudplatformv2_organization_accounts.accounts.sort_result
}