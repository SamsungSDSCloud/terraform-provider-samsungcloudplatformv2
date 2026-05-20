output "organization_units" {
  description = "List of Organization Units"
  value       = data.samsungcloudplatformv2_organization_units.org_units.organization_units
}