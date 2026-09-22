output "count" {
  description = "Total count"
  value       = data.samsungcloudplatformv2_organization_delegation_accounts.delegation_accounts.total_count
}

output "page" {
  description = "Page number"
  value       = data.samsungcloudplatformv2_organization_delegation_accounts.delegation_accounts.page
}

output "size" {
  description = "Page size"
  value       = data.samsungcloudplatformv2_organization_delegation_accounts.delegation_accounts.size
}

output "sort" {
  description = "Sort criteria"
  value       = data.samsungcloudplatformv2_organization_delegation_accounts.delegation_accounts.sort_result
}

output "delegation_accounts" {
  description = "Delegation Account list"
  value       = data.samsungcloudplatformv2_organization_delegation_accounts.delegation_accounts.delegation_accounts
}