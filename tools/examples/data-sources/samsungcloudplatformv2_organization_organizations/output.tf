output "organizations" {
  description = "Organizations"
  value       = data.samsungcloudplatformv2_organization_organizations.orgs.organizations
}

output "total_count" {
  description = "Total count"
  value       = data.samsungcloudplatformv2_organization_organizations.orgs.total_count
}

output "page" {
  description = "Page number"
  value       = data.samsungcloudplatformv2_organization_organizations.orgs.page
}

output "size" {
  description = "Page size"
  value       = data.samsungcloudplatformv2_organization_organizations.orgs.size
}

output "sort" {
  description = "Sort criteria"
  value       = data.samsungcloudplatformv2_organization_organizations.orgs.sort_result
}